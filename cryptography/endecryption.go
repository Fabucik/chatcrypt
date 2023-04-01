package cryptography

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

func Encrypt(pub rsa.PublicKey, message string) ([]byte, error) {
	encrypted, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		&pub,
		[]byte(message),
		nil,
	)

	if err != nil {
		return []byte(""), err
	}

	return encrypted, nil
}

func Decrypt(priv *rsa.PrivateKey, encrypted []byte) ([]byte, error) {
	decrypted, err := priv.Decrypt(
		rand.Reader,
		encrypted,
		&rsa.OAEPOptions{Hash: crypto.SHA256},
	)

	if err != nil {
		return []byte(""), err
	}

	return decrypted, nil
}
