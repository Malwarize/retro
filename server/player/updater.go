package player

import (
	"encoding/json"
	"github.com/Malwarize/retro/logger"
	"github.com/Malwarize/retro/shared"
	"io"
	"net/http"
	"os"
	"os/exec"
)

const ENDPOINT = "https://api.github.com/repos/Malwarize/retro/releases/latest"

type Updater struct {
	IsUpdateAvailable  bool
	RemoteVersion      string
	EnableUpdatePrompt bool
}

func NewUpdater() *Updater {
	return &Updater{
		EnableUpdatePrompt: true,
	}
}

type Release struct {
	Version string `json:"tag_name"`
}

func (u *Updater) GetRemoteVersion() (string, error) {
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

func (u *Updater) NeedsUpdate(currentVersion string) (bool, string) {
	remoteV, err := u.GetRemoteVersion()
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

func (u *Updater) Update() {
	go func() {
		if !u.IsUpdateAvailable {
			return
		}
		u.IsUpdateAvailable = false
		var DownloadEndpoint = "https://github.com/Malwarize/retro/releases/download/" + u.RemoteVersion + "/installer.tar.gz"
		logger.LogInfo("Downloading", DownloadEndpoint)
		req, err := http.Get(DownloadEndpoint)
		if err != nil {
			logger.LogError(
				logger.GError("Error Updating", err),
			)
			u.IsUpdateAvailable = true
			return
		}

		body, err := io.ReadAll(req.Body)
		if err != nil {
			logger.LogError(
				logger.GError("Error Updating", err),
			)
			u.IsUpdateAvailable = true
			return
		}
		defer req.Body.Close()
		logger.LogInfo("Downloaded", DownloadEndpoint)
		tmpFile, err := os.CreateTemp("", "retro")
		logger.LogInfo("Saving to temp file", tmpFile.Name())
		if err != nil {
			logger.LogError(
				logger.GError("Error Updating", err),
			)
			u.IsUpdateAvailable = true
			return
		}
		//defer os.Remove(tmpFile.Name())
		_, err = tmpFile.Write(body)
		err = tmpFile.Close()
		if err != nil {
			logger.LogError(
				logger.GError("Error Updating", err),
			)
			u.IsUpdateAvailable = true
			return
		}
		logger.LogInfo("Saved to temp file", tmpFile.Name())
		logger.LogInfo("Extracting", tmpFile.Name())
		cmd := exec.Command("tar", "-xf", tmpFile.Name(), "-C", "/tmp")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		logger.LogInfo("Extracting", tmpFile.Name())
		err = cmd.Run()
		if err != nil {
			logger.LogError(
				logger.GError("Error Updating", err),
			)
			u.IsUpdateAvailable = true
			return
		}
		logger.LogInfo("Running installer.sh")
		cmd = exec.Command("bash", "nohup /tmp/installer.sh &")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			logger.LogError(
				logger.GError("Error Updating", err),
			)
			u.IsUpdateAvailable = true
			return
		}
		logger.LogInfo("Cleaning up")
		logger.LogInfo("Update done")
		return
	}()
}

func (u *Updater) CheckUpdate() {
	if needUpdate, newVersion := u.NeedsUpdate(shared.Version); needUpdate {
		u.IsUpdateAvailable = true
		u.RemoteVersion = newVersion
	}
}
