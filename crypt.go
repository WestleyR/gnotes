// Created on: 2022-02-20

package gnotes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
)

func (c *cryptConfig) EncryptIfEnabled(data []byte) ([]byte, error) {
	if !c.Enable {
		log.Printf("Not encrypting notes.")
		return data, nil
	}

	if c.Key == "" {
		return []byte{}, fmt.Errorf("need a 16 bit key")
	}

	block, err := aes.NewCipher([]byte(c.Key))
	if err != nil {
		return []byte{}, fmt.Errorf("could not create new cipher: %s", err)
	}

	encData := make([]byte, aes.BlockSize+len(data))
	iv := encData[:aes.BlockSize]

	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return []byte{}, fmt.Errorf("could not encrypt: %s", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encData[aes.BlockSize:], data)

	log.Printf("Successfully encrypted data")

	return encData, nil
}

func (c *cryptConfig) DecryptIfEnabled(data []byte) ([]byte, error) {
	// Always try to decrypt the data, this way the user can disable encryption,
	// and get their encrypted notes, and it can be decrypted.

	if c.Key == "" {
		return []byte{}, fmt.Errorf("need a 16 bit key")
	}

	block, err := aes.NewCipher([]byte(c.Key))
	if err != nil {
		return []byte{}, fmt.Errorf("could not create new cipher: %s", err)
	}

	if len(data) < aes.BlockSize {
		return []byte{}, fmt.Errorf("invalid ciphertext block size")
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)

	return data, nil
}
