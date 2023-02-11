Quickstart
---
<!-- TOC -->
  * [Quickstart](#quickstart)
    * [Import](#import)
    * [Choose a Signing Method](#choose-a-signing-method)
    * [You can generate a key pair](#you-can-generate-a-key-pair)
    * [Or, load your own keys in and create a key pair](#or-load-your-own-keys-in-and-create-a-key-pair)
    * [Instantiate the manager, and you're off](#instantiate-the-manager-and-youre-off)
    * [Generate a token using the manager](#generate-a-token-using-the-manager)
    * [Generate a token with registered claims](#generate-a-token-with-registered-claims)
    * [Validate a token using the manager](#validate-a-token-using-the-manager)
<!-- TOC -->
### Import
```go
import "github.com/scottcagno/webslinger/pkg/web/jwt"
```

### Choose a Signing Method
```go
...
// RSA signing methods
_ = jwt.RS256
_ = jwt.RS384
_ = jwt.RS512

// HMAC signing methods
_ = jwt.HS256
_ = jwt.HS384
_ = jwt.HS512

// ECDSA signing methods
_ = jwt.ES256
_ = jwt.ES384
_ = jwt.ES512

// RSAPSS signing methods
_ = jwt.PS256
_ = jwt.PS384
_ = jwt.PS512
...
```

### You can generate a key pair
```go
...
// Generate key pair using RSA 512
keys := jwt.RS512.GenerateKeyPair()
...
```

### Or, load your own keys in and create a key pair
```go
...
// Read private key file
privateKeyData, err := os.ReadFile("mykeys/private.pem") 
if err != nil {
    log.Fatal(err)
}

// Read public key file
publicKeyData, err := os.ReadFile("mykeys/public.pem")
if err != nil {
    log.Fatal(err)
}

// Parse your private key
pri, err := jwt.ParsePrivateKeyFromPEM(privateKeyData)
if err != nil {
    log.Fatal(err)
}

// Parse your public key
pub, err := jwt.ParsePublicKeyFromPEM(publicKeyData)
if err != nil {
    log.Fatal(err)
}

// Create a *KeyPair to use with the manager
keys := &jwt.KeyPair{
    PrivateKey: pri,
    PublicKey:  pub,
}
...
```

### Instantiate the manager, and you're off
```go
...
// Instantiate the token manager
keys := jwt.HS256.GenerateKeyPair()
manager := jwt.NewTokenManager(jwt.HS256, keys)
...
```

### Generate a token using the manager
```go
...
// Generating a signed token with no claims
token, err := manager.GenerateToken(nil)
if err != nil {
    t.Fatal(err)
}
...
```

### Generate a token with registered claims
```go
...
// Set some time info
exp := time.Now().Add(1 * time.Hour).Unix()
now := time.Now().Unix()

// Create some registered claims
claims := &jwt.RegisteredClaims{
    Issuer:         "jon doe",
    Subject:        "your mom goes to college",
    Audience:       "anyone",
    ExpirationTime: NumericDate(exp),
    NotBeforeTime:  NumericDate(now),
    IssuedAtTime:   NumericDate(now),
    ID:             "25",
}

// Generating a signed token with registered claims
token, err := manager.GenerateToken(claims)
if err != nil {
    t.Fatal(err)
}
```


### Validate a token using the manager
```go
...
// Validate a token using the manager
validToken, err = manager.ValidateToken(token)
if err != nil {
    t.Fatal(err)
}
...
```


