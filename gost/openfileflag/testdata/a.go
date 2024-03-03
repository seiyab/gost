package testdata

import "os"

func _() {
	os.OpenFile("", os.O_RDONLY, 0600)
	os.OpenFile("", os.O_WRONLY|os.O_TRUNC, 0600)
	os.OpenFile("", os.O_WRONLY|os.O_APPEND, 0600)
	os.OpenFile("", os.O_WRONLY|os.O_EXCL, 0600)
	os.OpenFile("", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	os.OpenFile("", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	os.OpenFile("", os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)

	os.OpenFile("", os.O_WRONLY, 0600)             // want ".+"
	os.OpenFile("", os.O_WRONLY|os.O_CREATE, 0600) // want ".+"

	f1, err1 := os.OpenFile("", os.O_WRONLY, 0600) // want ".+"
	f1, err1 = os.OpenFile("", os.O_WRONLY, 0600)  // want ".+"
	_, err2 := os.OpenFile("", os.O_WRONLY, 0600)
	_, err2 = os.OpenFile("", os.O_WRONLY, 0600)

	markAsUsed(f1, err1, err2)
}

func markAsUsed(_ ...any) {}
