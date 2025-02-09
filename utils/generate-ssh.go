package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func GenerateSSHKey(email string) (string, string, error) {
	tinyID := fmt.Sprintf("%d", time.Now().UnixNano())[:8]
	privateKeyPath := filepath.Join(
		os.Getenv("HOME"),
		".ssh",
		"teagit_"+tinyID+"_ed25519",
	)
	publicKeyPath := privateKeyPath + ".pub"

	cmd := exec.Command(
		"ssh-keygen",
		"-t",
		"ed25519",
		"-C",
		email,
		"-f",
		privateKeyPath,
		"-N",
		"",
	)
	if err := cmd.Run(); err != nil {
		return "", "", fmt.Errorf("failed to generate SSH key: %v", err)
	}

	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", "", fmt.Errorf(
			"failed to read generated SSH private key: %v",
			err,
		)
	}
	publicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return "", "", fmt.Errorf(
			"failed to read generated SSH public key: %v",
			err,
		)
	}

	return string(privateKey), string(publicKey), nil
}
