package util

import (
	"bytes"
	"fmt"
	"os/exec"
)

// StartDockerContainer creates and runs a new Docker container
func StartDockerContainer(imageName string, containerName string) error {
	cmd := exec.Command("docker", "run", "-d", "--name", containerName, imageName)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		dockerError := stderr.String()
		return fmt.Errorf("failed to start Docker container: %w\nDocker error: %s", err, dockerError)
	}
	return nil
}
