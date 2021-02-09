package internal

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/google/tink/go/hybrid"
	"github.com/google/tink/go/insecurecleartextkeyset"
	"github.com/google/tink/go/keyset"
	"log"
	"os"
	"strings"
)

func Receive(numberOfSenders int) error {
	// Generate key pair
	privateKeyHandle, err := keyset.NewHandle(hybrid.ECIESHKDFAES128CTRHMACSHA256KeyTemplate())
	if err != nil {
		return fmt.Errorf("failed to generate transport key: %w", err)
	}

	// Obtain the public key
	publicKeyHandle, err := privateKeyHandle.Public()
	if err != nil {
		return fmt.Errorf("failed to generate public part of transport key: %w", err)
	}

	// Export the public key
	exportedPub := &keyset.MemReaderWriter{}
	if err = insecurecleartextkeyset.Write(publicKeyHandle, exportedPub); err != nil {
		return fmt.Errorf("failed to export public key: %w", err)
	}

	publicKey, err := proto.Marshal(exportedPub.Keyset)
	if err != nil {
		return fmt.Errorf("failed to marshal the public key: %w", err)
	}

	_, _ = fmt.Fprintln(os.Stderr, "Ask the sender(s) to execute the following command:")
	_, _ = fmt.Fprintln(os.Stderr, "")
	chunks("echo -n \"the secret\" | secret-sharer wrap --transport-key \""+base64.StdEncoding.EncodeToString(publicKey)+"\"", 80)
	_, _ = fmt.Fprintln(os.Stderr, "")

	for i := 0; i < numberOfSenders; i++ {
		_, _ = fmt.Fprintf(os.Stderr, "Input sender's data %d/%d:\n", i+1, numberOfSenders)

		// Get the input
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read sender's data: %w", err)
		}

		// Remove the newline
		input = strings.Replace(input, "\n", "", -1)

		// Decode the message
		_, plainText, err := decryptMessage(input, privateKeyHandle)
		if err != nil {
			return err
		}

		fmt.Println(plainText)
	}

	return nil
}

func decryptMessage(input string, privateKeyHandle *keyset.Handle) (string, string, error) {
	items := strings.Split(input, "$")
	if len(items) > 2 {
		return "", "", fmt.Errorf("ciphertext has too many parts: %d", len(items))
	}

	var contextInfo, cipherText string
	if len(items) == 2 {
		contextInfo, cipherText = items[0], items[1]
	} else {
		contextInfo, cipherText = "", items[0]
	}

	hd, err := hybrid.NewHybridDecrypt(privateKeyHandle)
	if err != nil {
		log.Fatal(err)
	}

	rawCipherText, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", "", fmt.Errorf("failed to decode the ciphertext: %w", err)
	}

	rawContextInfo, err := base64.StdEncoding.DecodeString(contextInfo)
	if err != nil {
		return "", "", fmt.Errorf("failed to decode the context info: %w", err)
	}

	plainText, err := hd.Decrypt(rawCipherText, rawContextInfo)
	if err != nil {
		return "", "", fmt.Errorf("failed to decrypt the ciphertext: %w", err)
	}

	return contextInfo, string(plainText), nil
}

func chunks(s string, chunkSize int) {
	if chunkSize >= len(s) {
		_, _ = fmt.Fprintln(os.Stderr, s)
		return
	}

	var chunks []string
	chunk := make([]rune, chunkSize)
	l := 0

	for _, r := range s {
		chunk[l] = r
		l++
		if l == chunkSize {
			chunks = append(chunks, string(chunk))
			l = 0
		}
	}

	if l > 0 {
		chunks = append(chunks, string(chunk[:l]))
	}

	for _, s := range chunks[:len(chunks)-1] {
		_, _ = fmt.Fprintln(os.Stderr, s+"\\")
	}

	if len(chunks) < 1 {
		_, _ = fmt.Fprintln(os.Stderr, s)
		return
	}

	//goland:noinspection GoNilness
	_, _ = fmt.Fprintln(os.Stderr, chunks[len(chunks)-1])
}

func Wrap(transportKey string, secret string, contextInfo string) error {
	rawTransportKey, err := base64.StdEncoding.DecodeString(transportKey)

	khPub, err := insecurecleartextkeyset.Read(keyset.NewBinaryReader(bytes.NewReader(rawTransportKey)))
	if err != nil {
		return fmt.Errorf("failed to read transport key %w", err)
	}

	he, err := hybrid.NewHybridEncrypt(khPub)
	if err != nil {
		log.Fatal(err)
	}

	ct, err := he.Encrypt([]byte(secret), []byte(contextInfo))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s$%s\n",
		base64.StdEncoding.EncodeToString([]byte(contextInfo)),
		base64.StdEncoding.EncodeToString(ct))

	return nil
}
