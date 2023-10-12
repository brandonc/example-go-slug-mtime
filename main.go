package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/hashicorp/go-slug"
)

func main() {
	buf := bytes.NewBuffer(nil)

	sourceDir := "./archive-dir"
	if len(os.Args) == 2 {
		sourceDir = os.Args[1]
	}

	meta, err := slug.Pack(sourceDir, buf, true)
	if err != nil {
		fmt.Printf("Unexpected error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Packed %q: %d files / %d bytes\n", sourceDir, len(meta.Files), meta.Size)

	targetDir, err := os.MkdirTemp("", "archive-dir")
	if err != nil {
		fmt.Printf("Unexpected error creating destination directory: %v\n", err)
		os.Exit(1)
	}

	// Unpacking a slug is done by calling the Unpack function with an
	// io.Reader to read the slug from and a directory path of an existing
	// directory to store the unpacked configuration files.
	if err := slug.Unpack(buf, targetDir); err != nil {
		fmt.Printf("Unexpected error unpacking: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(targetDir)

	a, err := hashPolicies(targetDir)
	if err != nil {
		fmt.Printf("Unexpected error hashing unpacked directory: %v\n", err)
		os.Exit(1)
	}

	b, err := hashPolicies(sourceDir)
	if err != nil {
		fmt.Printf("Unexpected error hashing source directory: %v\n", err)
		os.Exit(1)
	}

	if a != b {
		fmt.Println("❌ directory hashes are NOT equal after unpack.")
		os.Exit(1)
	}

	fmt.Println("✅ directory hashes are equal after unpacking")
}

func hashPolicies(path string) (string, error) {
	body := bytes.NewBuffer(nil)
	file, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if !file.Mode().IsDir() {
		return "", fmt.Errorf("the path is not a directory")
	}

	_, err = slug.Pack(path, body, true)
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	hash.Write(body.Bytes())
	chksum := hex.EncodeToString(hash.Sum(nil))

	return chksum, nil
}
