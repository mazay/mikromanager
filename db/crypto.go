package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

// DecryptString takes a string that is encrypted with the given key and returns
// the decrypted string. The encrypted string is expected to be the output of
// EncryptString and the key is expected to be the same as the one used for
// encryption. If the key is incorrect or the encrypted string is tampered with,
// an error is returned.
func DecryptString(cryptoText string, keyString string) (plainTextString string, err error) {
	encrypted, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}
	if len(encrypted) < aes.BlockSize {
		return "", fmt.Errorf("cipherText too short. It decodes to %v bytes but the minimum length is 16", len(encrypted))
	}

	decrypted, err := decryptAES(hashTo32Bytes(keyString), encrypted)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

// decryptAES decrypts the given data using the given key with AES encryption
// using a cipher feedback mode of operation. The first 16 bytes of the data
// are expected to be the initialization vector. The function returns an error
// if the decryption fails.
func decryptAES(key, data []byte) ([]byte, error) {
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(data, data)
	return data, nil
}

// hashTo32Bytes takes a string, hashes it with SHA-256 and returns the first
// 32 bytes of the hash. This is used to generate a key for AES encryption
// from a given string.
func hashTo32Bytes(input string) []byte {
	data := sha256.Sum256([]byte(input))
	return data[0:]

}

// EncryptString encrypts the given plaintext string using the given keyString
// with AES encryption and returns the encrypted data as a URL-safe base64
// encoded string. The key is first hashed with SHA-256 and the first 32 bytes
// of the hash are used as the key for AES encryption. The function returns an
// error if the encryption fails.
func EncryptString(plainText string, keyString string) (cipherTextString string, err error) {
	key := hashTo32Bytes(keyString)
	encrypted, err := encryptAES(key, []byte(plainText))
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(encrypted), nil
}

// encryptAES encrypts the given data using the given key with AES encryption
// in cipher feedback mode. A random initialization vector (IV) is generated
// and prepended to the encrypted data. The function returns the IV combined
// with the encrypted data or an error if the encryption fails.
func encryptAES(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// create two 'windows' in to the output slice.
	output := make([]byte, aes.BlockSize+len(data))
	iv := output[:aes.BlockSize]
	encrypted := output[aes.BlockSize:]

	// populate the IV slice with random data.
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)

	// note that encrypted is still a window in to the output slice
	stream.XORKeyStream(encrypted, data)
	return output, nil
}
