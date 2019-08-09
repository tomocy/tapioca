package infra

import (
	"encoding/json"
	"io"
	"time"
)

func readJSON(src io.Reader, dest interface{}) error {
	return json.NewDecoder(src).Decode(dest)
}

func today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}
