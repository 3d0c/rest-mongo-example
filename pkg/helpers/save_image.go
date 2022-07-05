package helpers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/teal-seagull/lyre-be-v4/pkg/config"
)

// ParseAndSaveImage decodes base64 payload and saves file into config.Static folder
func ParseAndSaveImage(p string) (string, error) {
	var (
		buff     bytes.Buffer = bytes.Buffer{}
		err      error
		idx      int
		filename string
		dec      io.Reader
	)

	init := strings.Index(p, "/")
	end := strings.Index(p, ";")
	filetype := p[init+1 : end]

	filename = uuid.New().String() + "." + filetype

	if idx = strings.Index(p, ","); idx < 0 {
		return "", fmt.Errorf("error parsing base64 payload")
	}

	dec = base64.NewDecoder(base64.StdEncoding, strings.NewReader(p[idx+1:]))

	if _, err = buff.ReadFrom(dec); err != nil {
		return "", err
	}

	if err = os.MkdirAll(config.TheConfig().Server.Static, 0755); err != nil {
		return "", err
	}

	filename = filepath.Join(config.TheConfig().Server.Static, filename)

	if err = ioutil.WriteFile(filename, buff.Bytes(), 0600); err != nil {
		return "", err
	}

	return filename, nil
}
