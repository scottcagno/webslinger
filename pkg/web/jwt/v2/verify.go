package v2

type VerifiedToken struct {
	Token     []byte
	Header    any
	Payload   any
	Signature any
	Claims    any
}

func Verify(alg Alg, k PublicKey, token []byte) (*VerifiedToken, error) {
	if len(token) == 0 {
		return nil, ErrMissing
	}
	header, payload, signature, err := decode()
}
