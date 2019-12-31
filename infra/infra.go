package infra

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
)

func createWorkspace() error {
	name := configFilename()
	if _, err := os.Stat(name); err == nil {
		return nil
	}

	dir := workspaceName()
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	f, err := os.Create(name)
	if err != nil {
		return err
	}

	return f.Close()
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
	return filepath.Join(os.Getenv("HOME"), ".tapioca")
}

func readJSON(src io.Reader, dest interface{}) error {
	return json.NewDecoder(src).Decode(dest)
}

func writeJSON(dest io.Writer, src interface{}) error {
	return json.NewEncoder(dest).Encode(src)
}

type config struct {
	GitHub oauth2Config
}

type oauth2Config struct {
	AccessToken *oauth2.Token
}
