package testdata

import (
	"os"
	"path"
	"path/filepath"
)

func _() {
	os.RemoveAll(path.Join("a", "b")) // want ".+"
	os.RemoveAll(filepath.Join("a", "b"))
}
