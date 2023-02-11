###
<div style="text-align: center">
    <img src="https://jwt.io/img/pic_logo.svg" alt="jwt image">
</div>
<h2 style="text-align:center">JSON Web Tokens</h2>

---

### What is JSON Web Token?
JSON Web Token (JWT) is an open standard [RFC 7519](https://tools.ietf.org/html/rfc7519) 
that defines a compact and self-contained way for securely transmitting information 
between parties as a JSON object. This information can be verified and trusted because 
it is digitally signed. JWTs can be signed using a secret (with the HMAC algorithm) or a 
public/private key pair using RSA or ECDSA.

Although JWTs can be encrypted to also provide secrecy between parties, we will focus on 
signed tokens. Signed tokens can verify the integrity of the claims contained within it, 
while encrypted tokens hide those claims from other parties. When tokens are signed using 
public/private key pairs, the signature also certifies that only the party holding the 
private key is the one that signed it.
---

### What is the JSON Web Token Structure?
In its compact form, JSON Web Tokens consist of three parts separated by dots `(.)`, 
which are:
- Header
- Payload
- Signature

Therefore, a JWT typically looks like the following.

`xxxx.yyyy.zzzz`

Let's break down the different parts.

---

#### Header
The header *typically* consists of two parts: the type of the token, which is JWT,
and the signing algorithm being used, such as HMAC, SHA256 or RSA.

For example:

```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```
Then, this JSON is **Base64Url** encoded to form the ***first*** part of the JWT.

---

#### Payload
The second part of the token is the payload, which contains the claims. Claims 
are statements about an entity (typically, the user) and additional data. There 
are three types of claims: registered, public, and private claims.

- **Registered claims:** These are a set of predefined claims which are not 
  mandatory but recommended, to provide a set of useful, interoperable claims. 
  Some of them are: iss (issuer), exp (expiration time), sub (subject), aud 
  (audience), and [others](https://tools.ietf.org/html/rfc7519#section-4.1).


- **Public claims:** These can be defined at will by those using JWTs. But to
  avoid collisions they should be defined in the IANA JSON Web Token Registry 
  or be defined as a URI that contains a collision resistant namespace.


- **Private claims:** These are the custom claims created to share information 
  between parties that agree on using them and are neither registered nor, public.

---
<p style="text-align: center">
    Do note that the claim names are only three characters<br>
    in length as JWT is meant to be compact.
</p>

---

An example payload could be:

```json
{
  "sub": "1234567890",
  "name": "John Doe",
  "admin": true
}
```

Then, this JSON is **Base64Url** encoded to form the ***second*** part of the JWT.

---
<p style="text-align: center">
    Do note that for signed tokens this information, though protected against <br>
    tampering, is readable by anyone. Do not put secret information in the <br>
    payload or header elements of a JWT unless it is encrypted.
</p>

---

#### Signature
To create the signature part you have to take the encoded header, the encoded payload, 
a secret, the algorithm specified in the header, and sign that.

For example if you want to use the HMAC SHA256 algorithm, the signature will be created 
in the following way:

```javascript
HMACSHA256(
  base64UrlEncode(header) + "." +
  base64UrlEncode(payload),
  secret)
```

The signature is used to verify the message wasn't changed along the way, and, in the 
case of tokens signed with a private key, it can also verify that the sender of the JWT
is who it says it is.

---

### Putting it all together
The output is three Base64-URL strings separated by dots that can be easily passed in 
HTML and HTTP environments, while being more compact when compared to XML-based standards
such as SAML.

The following shows a JWT that has the previous header and payload encoded, and it is 
signed with a secret.

<div style="text-align: center">
<img alt="jwt encoded image" src="https://cdn.auth0.com/content/jwt/encoded-jwt3.png" width="500"/>
</div>

If you want to play with JWT and put these concepts into practice, you can use 
[jwt.io Debugger](https://jwt.io/#debugger-io) to decode, verify, and generate JWTs.

For more information, please check out the JWT page I totally ripped off at
[https://jwt.io/introduction](https://jwt.io/introduction)