package cryptography

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
)

func GenerateKeyPair() (*rsa.PrivateKey, rsa.PublicKey) {
	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	pubKey := privKey.PublicKey

	return privKey, pubKey
}

func SaveKeysToFiles(privKey *rsa.PrivateKey, pubKey rsa.PublicKey, savePath string) error {

	pemPriv := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privKey),
		},
	)

	pemPub := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(&pubKey),
		},
	)

	if savePath == "" {
		ex, _ := os.Executable()
		savePath = filepath.Dir(ex) + "/keys"
	}

	err := os.WriteFile(savePath+"/chatcrypt.priv", pemPriv, 0600)
	if err != nil {
		return err
	}

	err = os.WriteFile(savePath+"/chatcrypt.pub", pemPub, 0600)
	if err != nil {
		return err
	}

	return nil
}

func ReadKeysFromFiles(keyPath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	if keyPath == "" {
		ex, _ := os.Executable()
		keyPath = filepath.Dir(ex) + "/keys"
	}

	privBytes, err := os.ReadFile(keyPath + "/chatcrypt.priv")
	if err != nil {
		return nil, nil, err
	}

	pubBytes, err := os.ReadFile(keyPath + "/chatcrypt.pub")
	if err != nil {
		return nil, nil, err
	}

	privBlock, _ := pem.Decode(privBytes)
	pubBlock, _ := pem.Decode(pubBytes)

	privKey, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	pubKey, err := x509.ParsePKCS1PublicKey(pubBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	return privKey, pubKey, nil
}
