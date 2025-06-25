package main

import (
	"archive/zip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"os"
)

// Create a ZIP archive from a list of files
func createZip(output string, files []string) error {
	newZipFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	for _, file := range files {
		err = addFileToZip(zipWriter, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func addFileToZip(zipWriter *zip.Writer, filename string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = filename

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, fileToZip)
	return err
}

// Encrypt inputFile and save to outputFile using AES-GCM
func encryptFile(key []byte, inputFile, outputFile string) error {
	in, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, in, nil)

	err = os.WriteFile(outputFile, ciphertext, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Decrypt inputFile and save to outputFile
func decryptFile(key []byte, inputFile, outputFile string) error {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	return os.WriteFile(outputFile, plaintext, 0644)
}

func main() {
	keyPtr := flag.String("key", "", "Encryption key (32 bytes)")
	outputPtr := flag.String("out", "backup.zip.aes", "Output file name")
	decryptPtr := flag.Bool("decrypt", false, "Decryption mode (decrypt file)")
	flag.Parse()
	args := flag.Args()

	if len(*keyPtr) != 32 {
		fmt.Println("Error: key must be exactly 32 bytes")
		return
	}

	if *decryptPtr {
		// Decryption mode: expect 1 input file
		if len(args) != 1 {
			fmt.Println("Error: specify exactly one encrypted file to decrypt")
			return
		}
		inputFile := args[0]
		outputFile := *outputPtr

		err := decryptFile([]byte(*keyPtr), inputFile, outputFile)
		if err != nil {
			fmt.Println("Decryption error:", err)
			return
		}
		fmt.Println("File successfully decrypted and saved as:", outputFile)

	} else {
		// Archive + encrypt mode
		if len(args) == 0 {
			fmt.Println("Error: specify files to archive")
			return
		}
		tempZip := "temp_backup.zip"

		err := createZip(tempZip, args)
		if err != nil {
			fmt.Println("Archiving error:", err)
			return
		}

		err = encryptFile([]byte(*keyPtr), tempZip, *outputPtr)
		if err != nil {
			fmt.Println("Encryption error:", err)
			return
		}

		os.Remove(tempZip)
		fmt.Println("Backup successfully created and encrypted in file:", *outputPtr)
	}
}
