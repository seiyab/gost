package testdata

import "github.com/pkg/errors"

func _() {
	var err error

	errors.WithStack(err)
	errors.WithMessage(err, "")
	errors.WithMessagef(err, "")
	errors.Wrap(err, "")
	errors.Wrapf(err, "")

	someFunc(nil)

	errors.WithStack(nil)        // want ".+"
	errors.WithMessage(nil, "")  // want ".+"
	errors.WithMessagef(nil, "") // want ".+"
	errors.Wrap(nil, "")         // want ".+"
	errors.Wrapf(nil, "")        // want ".+"
}

func someFunc(e error) error {
	return e
}
