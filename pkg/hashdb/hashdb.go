package hashdb

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
)

// bin is our shorthand for how we encode binary data
var bin = binary.BigEndian

// HashFn defines our hash function type. This is here
// in order to make it easier to swap out at any point.
type HashFn func(msg []byte) uint32

// cbdHash is just our default hash implementation for now.
func cbdHash(b []byte) uint32 {
	j := uint32(5381)
	size := len(b)
	for size != 0 {
		size--
		j += j << 5
		j ^= uint32(b[size])
	}
	j &= 0xffffffff
	return j
}

const recordHdr = 8

// record is a data entry record
type record struct {
	key []byte
	val []byte
}

func (r *record) String() string {
	return fmt.Sprintf("[%.4d,%.4d,%q,%q]", len(r.key), len(r.val), r.key, r.val)
}

func decRecord(data []byte) *record {
	// get the 4 byte length of the record key
	klen := bin.Uint32(data[0:4])
	// get the 4 byte length of the record value
	vlen := bin.Uint32(data[4:8])
	// initialize our new record
	rec := &record{
		key: make([]byte, klen),
		val: make([]byte, vlen),
	}
	// copy the key and value data and return
	copy(rec.key, data[recordHdr:recordHdr+klen])
	copy(rec.val, data[recordHdr+klen:recordHdr+klen+vlen])
	// return record
	return rec
}

// {0x00,0x00,0x00,0x00,0xff,0xff,0xff,0xff,0x00,0x00,0x00,0x00,0xff,0xff,0xff}
//   0    1    2    3    4    5    6    7    8    9    10   11   12   13,  14

func encRecord(k, v []byte) ([]byte, int) {
	// record header length
	off := recordHdr
	// create our record
	rec := make([]byte, off+len(k)+len(v))
	// 4 byte length of the record key
	bin.PutUint32(rec[0:4], uint32(len(k)))
	// 4 byte length of the record value
	bin.PutUint32(rec[4:8], uint32(len(v)))
	// a variable length binary encoded key
	off += copy(rec[off:], k)
	// a variable length binary encoded value
	off += copy(rec[off:], v)
	// return record
	return rec, off
}

func toString(rec []byte) string {
	return decRecord(rec).String()
}

type index struct {
	hash   uint32
	offset uint32
}

type table struct {
	hasher HashFn
	data   io.ReadWriteSeeker
	buf    []byte
	fp     *os.File
}

func (t *table) readU32At(offset uint32) (uint32, error) {
	curr, err := t.data.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return 0, err
	}
	defer func(data io.ReadWriteSeeker, offset int64, whence int) {
		_, err = data.Seek(offset, whence)
		if err != nil {
			panic(err)
		}
	}(t.data, curr, io.SeekStart)
	_, err = t.data.Read(t.buf)
	if err != nil {
		return 0, err
	}
	return bin.Uint32(t.buf), nil
}

func (t *table) writeU32At(u32, offset uint32) error {
	curr, err := t.data.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return err
	}
	defer func(data io.ReadWriteSeeker, offset int64, whence int) {
		_, err = data.Seek(offset, whence)
		if err != nil {
			panic(err)
		}
	}(t.data, curr, io.SeekStart)
	bin.PutUint32(t.buf, u32)
	_, err = t.data.Write(t.buf)
	if err != nil {
		return err
	}
	return nil
}

func (t *table) find(k []byte) (uint32, error) {
	var err error
	var h0, h1, h1e, last uint32
	// compute the hash value (H) of the key
	h := t.hasher(k)
	// read the hash table offset (O) at the location: 4 + (H & 0xff) * 4
	h0 = 4 + (h&0xff)*4
	h1, err = t.readU32At(h0)
	if err != nil {
		panic("bad read")
	}
	h1e, err = t.readU32At(h0 + 4)
	// read the next hash table offset (O2)--the number of slots in this hash table
	// is (O2 - O) / 8
	last = (h1e - h1) / 8

	return last, nil
}

func readKeyOff(keyHash uint32, rd io.ReaderAt) uint32 {
	buf := make([]byte, 4)
	_, err := rd.ReadAt(buf, int64(keyHash))
	if err != nil {
		return math.MaxUint32
	}
	return bin.Uint32(buf)
}

func writeKey(hasher HashFn, wr io.WriterAt, k []byte) {
	kh := getKeyHash(hasher, k)
	_ = kh
}

func getKeyHash(hasher HashFn, k []byte) uint32 {
	// h := hasher(k)
	// return 4 + (h&0xff)*4
	return hasher(k) & 0xff
}

type hashdb struct {
	index  *os.File
	data   *os.File
	buf    [4]byte
	hasher HashFn
}

func newHashDB(hasher HashFn, index, data *os.File) *hashdb {
	return &hashdb{
		index:  index,
		data:   data,
		buf:    [4]byte{},
		hasher: hasher,
	}
}

func (db *hashdb) get(k []byte) ([]byte, error) {
	// calculate key hash
	h := db.hasher(k) & 0xff
	// read the data at the key hash offset
	_, err := db.index.ReadAt(db.buf[:], int64(h))
	if err != nil {
		return nil, err
	}
	var off, size uint32
	// off is the data offset for key k
	off = bin.Uint32(db.buf[:])
	// read the size of our data
	_, err = db.data.ReadAt(db.buf[:], int64(off))
	if err != nil {
		return nil, err
	}
	size = bin.Uint32(db.buf[:])
	// create a byte slice to hold our data then read
	data := make([]byte, size)
	_, err = db.data.ReadAt(data, int64(off))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (db *hashdb) put(k, v []byte) error {
	// calculate key hash
	h := db.hasher(k) & 0xff
	// read the data at the key hash offset to make sure
	// the slot is available to write into later on
	_, err := db.index.ReadAt(db.buf[:], int64(h))
	if err != nil {
		return err
	}
	var off uint32
	var size int
	// off is the data offset for key k
	off = bin.Uint32(db.buf[:])
	if off == 0 {
		// there is not an existing record, treat
		// this like an insert.
		data := make([]byte, 12+len(k)+len(v))
		// encode the total record size
		bin.PutUint32(data[0:4], uint32(len(data)))
		// encode the key length
		bin.PutUint32(data[4:8], uint32(len(k)))
		// encode the val length
		bin.PutUint32(data[8:12], uint32(len(v)))
		// write the actual key and value
		copy(data[12:], k)
		copy(data[12+len(k):], v)
		size, err = db.data.Write(data)
		if err != nil {
			return err
		}
		// calculate the offset

	}

	_ = size

	return nil
}
