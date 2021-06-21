package file

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// Checksum calculates file checksum given its path.
func Checksum(filePath string) (checksum string, err error) {
	hasher := sha256.New()
	f, err := os.Open(filePath)
	if err != nil {
		err = fmt.Errorf("failed reading image path: %v", err)
		return
	}
	defer f.Close()
	if _, err = io.Copy(hasher, f); err != nil {
		err = fmt.Errorf("failed copying image content to hash: %v", err)
		return
	}

	checksum = hex.EncodeToString(hasher.Sum(nil))
	return
}
