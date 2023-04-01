package networking

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Fabucik/chatcrypt/cryptography"
	"github.com/Fabucik/chatcrypt/entities"
)

func SendMessage(message string, receiverPubKey rsa.PublicKey, senderPubKey rsa.PublicKey, senderPrivKey *rsa.PrivateKey, host string, port int) error {
	encryptedMessage, err := cryptography.Encrypt(receiverPubKey, message)
	if err != nil {
		return err
	}

	messageSignature, err := cryptography.SignMessage(senderPrivKey, encryptedMessage)
	if err != nil {
		return err
	}

	postBody, _ := json.Marshal(map[string][]byte{
		"to":        x509.MarshalPKCS1PublicKey(&receiverPubKey),
		"from":      x509.MarshalPKCS1PublicKey(&senderPubKey),
		"message":   encryptedMessage,
		"signature": messageSignature,
	})

	bufferedBody := bytes.NewBuffer(postBody)

	res, err := http.Post("http://"+host+":"+strconv.Itoa(port)+"/sendmessage", "application/json", bufferedBody)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		errorText := "Got bad status code: " + strconv.Itoa(res.StatusCode)
		return errors.New(errorText)
	}

	return nil
}

func FetchUnreadMessages(pubKey rsa.PublicKey, privKey *rsa.PrivateKey, host string, port int) (entities.AllMessages, error) {
	message := "AUTHENTICITY VERIFICATION"
	signature, err := cryptography.SignMessage(privKey, []byte(message))
	if err != nil {
		return entities.AllMessages{}, err
	}

	postBody, _ := json.Marshal(map[string][]byte{
		"publicKey": x509.MarshalPKCS1PublicKey(&pubKey),
		"signature": signature,
	})

	bufferedBody := bytes.NewBuffer(postBody)

	res, err := http.Post("http://"+host+":"+strconv.Itoa(port)+"/fetchmessages", "application/json", bufferedBody)
	if err != nil {
		return entities.AllMessages{}, err
	}

	if res.StatusCode != 200 {
		errorText := "Got bad status code: " + strconv.Itoa(res.StatusCode)
		return entities.AllMessages{}, errors.New(errorText)
	}

	defer res.Body.Close()

	var allMessages entities.AllMessages
	err = json.NewDecoder(res.Body).Decode(&allMessages)
	if err != nil {
		return entities.AllMessages{}, err
	}

	return allMessages, nil
}
