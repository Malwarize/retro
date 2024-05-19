package updater

import (
	"encoding/json"
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

func Update(version string) error {
	var DownloadEndpoint = "https://github.com/Malwarize/retro/releases/download/" + version + "/installer.tar.gz"

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
	cmd := exec.Command("tar", "-xvf", tmpFile.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	cmd = exec.Command("bash", "installer.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
