package infra

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/oauth2"
)

func today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func loadConfig() (*config, error) {
	name := configFilename()
	src, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	var loaded *config
	if err := readJSON(src, &loaded); err != nil {
		return nil, err
	}

	return loaded, nil
}

func saveConfig(cnf *config) error {
	name := configFilename()
	dest, err := os.OpenFile(name, os.O_WRONLY, 0700)
	if err != nil {
		return err
	}
	defer dest.Close()

	return writeJSON(dest, cnf)
}

func configFilename() string {
	return filepath.Join(workspaceName(), "config.json")
}

func workspaceName() string {
	return filepath.Join(os.Getenv("HOME"), "./tapioca")
}

func readJSON(src io.Reader, dest interface{}) error {
	return json.NewDecoder(src).Decode(dest)
}

func writeJSON(dest io.Writer, src interface{}) error {
	return json.NewEncoder(dest).Encode(src)
}

type config struct {
	github githubConfig
}

type githubConfig struct {
	AccessToken *oauth2.Token
}
