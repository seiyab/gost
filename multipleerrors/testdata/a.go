package testdata

import "errors"

func _() {
	var es []error
	var e error
	var errs error

	errs = errors.Join(errs, e)
	x := errors.Join(errs, e)
	y := errors.Join(es...)

	errors.Join(errs, e)      // want ".+"
	errs = errors.Join(e, e)  // want ".+"
	errs = errors.Join(errs)  // want "."
	errs = errors.Join(es...) // want ".+"

	markAsUsed(x, y)
}

func markAsUsed(_ ...any) {}
