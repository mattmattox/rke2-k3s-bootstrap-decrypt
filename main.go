package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/pbkdf2"
	"sigs.k8s.io/yaml"
)

type DecryptedData struct {
	ClientCAKey   string `json:"ClientCAKey"`
	ETCDPeerCA    string `json:"ETCDPeerCA"`
	ETCDPeerCAKey string `json:"ETCDPeerCAKey"`
	ETCDServerCA  string `json:"ETCDServerCA"`
}

const Version = "1.0.0"

func main() {
	// Define a flag for the passphrase
	passphrase := flag.String("passphrase", "", "The passphrase to decrypt the data.")
	bootstrapFile := flag.String("bootstrap-file", "bootstrap", "The name of the bootstrap file to read encrypted data from.")
	outputYAMLFile := flag.String("output-yaml-file", "decrypted.yaml", "The name of the YAML file to write decrypted data to.")
	flag.Parse()

	// Check if the passphrase is provided
	if *passphrase == "" {
		fmt.Println("Usage: main -passphrase=<passphrase> [options]")
		fmt.Println("Note: The token is stored in the format \"K10<CA-HASH>::<USERNAME>:<PASSWORD>\". Provide only the password at the end.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Read the encrypted data from the file.
	encryptedData, err := os.ReadFile(*bootstrapFile)
	if err != nil {
		log.Fatal("Error reading encrypted data from file:", err)
	}

	// Decrypt data
	decryptedData, err := decrypt(*passphrase, []byte(encryptedData))
	if err != nil {
		log.Fatal("Decryption error:", err)
	}

	// Write decrypted data to YAML file
	yamlData, err := yaml.JSONToYAML(decryptedData)
	if err != nil {
		log.Fatal("Error converting JSON to YAML:", err)
	}

	err = writeYAMLToFile(*outputYAMLFile, yamlData)
	if err != nil {
		log.Fatal("Error writing decrypted data to YAML file:", err)
	}
}

func decrypt(passphrase string, ciphertext []byte) ([]byte, error) {
	parts := strings.SplitN(string(ciphertext), ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid cipher text, not : delimited")
	}

	clearKey := pbkdf2.Key([]byte(passphrase), []byte(parts[0]), 4096, 32, sha1.New)
	key, err := aes.NewCipher(clearKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(key)
	if err != nil {
		return nil, err
	}

	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	return gcm.Open(nil, data[:gcm.NonceSize()], data[gcm.NonceSize():], nil)
}

// writeYAMLToFile writes YAML data to a file
func writeYAMLToFile(fileName string, yamlData []byte) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(yamlData)
	if err != nil {
		return err
	}

	return nil
}
