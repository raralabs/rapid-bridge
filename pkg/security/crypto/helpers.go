package hybridcrypto

import "strings"

func SplitMessage(message string) []string {
	return strings.Split(message, "-")
}

func CreateMessageToSign(ciphertext, aesKey, nonce []byte) []byte {
	messageToSign := make([]byte, 0, len(ciphertext)+1+len(aesKey)+1+len(nonce))
	messageToSign = append(messageToSign, ciphertext...)
	messageToSign = append(messageToSign, '-')
	messageToSign = append(messageToSign, aesKey...)
	messageToSign = append(messageToSign, '-')
	messageToSign = append(messageToSign, nonce...)

	return messageToSign
}
