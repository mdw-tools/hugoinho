package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mdwhatcott/static/contracts"
)

func LoadContent(folder string) map[contracts.Path]contracts.File {
	content := make(map[contracts.Path]contracts.File)
	_ = filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		data, _ := ioutil.ReadFile(path)
		content[contracts.Path(strings.TrimPrefix(path, folder))] = contracts.File(data)
		return nil
	})
	return content
}
