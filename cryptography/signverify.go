package cryptography

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

func SignMessage(privKey *rsa.PrivateKey, message []byte) ([]byte, error) {
	messageHash := sha256.New()
	_, err := messageHash.Write(message)
	if err != nil {
		return nil, err
	}

	messageHashSum := messageHash.Sum(nil)

	signature, err := rsa.SignPSS(rand.Reader, privKey, crypto.SHA256, messageHashSum, nil)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

func IsVerificationSuccessful(pubKey rsa.PublicKey, messageSum []byte, signature []byte) (bool, error) {
	err := rsa.VerifyPSS(&pubKey, crypto.SHA256, messageSum, signature, nil)
	if err != nil {
		return false, err
	}

	return true, nil
}
