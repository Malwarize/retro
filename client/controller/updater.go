package controller

import (
	"encoding/json"
	"fmt"
	"github.com/Malwarize/retro/logger"
	"github.com/Malwarize/retro/shared"
	"io"
	"net/http"
	"os"
	"os/exec"
)

const ENDPOINT = "https://api.github.com/repos/Malwarize/retro/releases/latest"

type Release struct {
	Version string `json:"tag_name"`
}

func GetRemoteVersion() (string, error) {
	req, err := http.Get(ENDPOINT)

	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}

	release := Release{}
	err = json.Unmarshal(body, &release)
	if err != nil {
		return "", err
	}
	return release.Version, nil

}

func NeedsUpdate(currentVersion string) (bool, string) {
	remoteV, err := GetRemoteVersion()
	if currentVersion != remoteV && err == nil {
		return true, remoteV
	}
	return false, ""
}

type asset struct {
	Url string `json:"browser_download_url"`
}

type Download struct {
	Assets []asset `json:"assets"`
}

func Update() error {
	needsUpdate, newVersion := NeedsUpdate(shared.Version)
	if !needsUpdate {
		return fmt.Errorf("No update available")
	}

	var DownloadEndpoint = "https://github.com/Malwarize/retro/releases/download/" + newVersion + "/installer.tar.gz"
	logger.LogInfo("Downloading", DownloadEndpoint)
	req, err := http.Get(DownloadEndpoint)
	if err != nil {
		return logger.LogError(
			logger.GError("Error Updating", err),
		)

	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return logger.LogError(
			logger.GError("Error Updating", err),
		)
	}
	defer req.Body.Close()
	logger.LogInfo("Downloaded", DownloadEndpoint)
	tmpFile, err := os.CreateTemp("", "retro")
	logger.LogInfo("Saving to temp file", tmpFile.Name())
	if err != nil {
		return logger.LogError(
			logger.GError("Error Updating", err),
		)

	}
	//defer os.Remove(tmpFile.Name())
	_, err = tmpFile.Write(body)
	err = tmpFile.Close()
	if err != nil {
		return logger.LogError(
			logger.GError("Error Updating", err),
		)
	}
	logger.LogInfo("Saved to temp file", tmpFile.Name())
	logger.LogInfo("Extracting", tmpFile.Name())
	cmd := exec.Command("tar", "-xf", tmpFile.Name(), "-C", "/tmp")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	logger.LogInfo("Extracting", tmpFile.Name())
	err = cmd.Run()
	if err != nil {
		return logger.LogError(
			logger.GError("Error Updating", err),
		)
	}
	logger.LogInfo("Running installer.sh")
	os.Chmod("/tmp/installer.sh", 0777)
	cmd = exec.Command("bash", "/tmp/installer.sh")
	if err := cmd.Start(); err != nil {
		return logger.LogError(
			logger.GError("Error Updating", err),
		)
	}
	if err := cmd.Process.Release(); err != nil {
		return logger.LogError(logger.GError("Error releasing process:", err))

	}
	logger.LogInfo("Cleaning up")
	logger.LogInfo("Update done")
	return nil

}
