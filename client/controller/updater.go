package controller

import (
	"encoding/json"
	"fmt"
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
		return fmt.Errorf("⛔ No update available")
	}

	var DownloadEndpoint = "https://github.com/Malwarize/retro/releases/download/" + newVersion + "/installer.tar.gz"
	fmt.Println("⬇️ Downloading", DownloadEndpoint)
	req, err := http.Get(DownloadEndpoint)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	tmpFile, err := os.CreateTemp("", "retro")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())
	_, err = tmpFile.Write(body)
	err = tmpFile.Close()
	if err != nil {
		return err
	}
	fmt.Println("💿  Saved to temp file", tmpFile.Name())
	cmd := exec.Command("tar", "-xf", tmpFile.Name(), "-C", "/tmp")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("📁 Extracting", tmpFile.Name())
	err = cmd.Run()
	if err != nil {
		return err
	}
	fmt.Println("🚀 Running installer.sh")
	os.Chmod("/tmp/installer.sh", 0777)
	cmd = exec.Command("bash", "/tmp/installer.sh")
	if err := cmd.Run(); err != nil {
		return err
	}
	fmt.Println("✅ Update done")
	return nil

}
