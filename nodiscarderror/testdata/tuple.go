package testdata

import "os"

func _() (startError error, endError error) {
	f, err := os.Open("file")
	if err != nil {
		return err, nil
	}

	if _, err := f.Read(nil); err != nil {
		return nil, nil // want ".+"
	}

	if err := f.Close(); err != nil {
		return nil, err
	}

	return nil, nil
}
