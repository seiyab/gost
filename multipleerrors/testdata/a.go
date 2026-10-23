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
	errs = errors.Join(e, e)  // want "original error is discarded"
	errs = errors.Join(errs)  // want "."
	errs = errors.Join(es...) // want "original error is discarded"

	markAsUsed(x, y)
}

func markAsUsed(_ ...any) {}

func reuseAfterHandling(run func() error, collected []error) error {
	err := run()
	if err != nil {
		return err
	}
	err = errors.Join(collected...)
	return err
}

func reuseWithExplicitArguments(run func() error, first, second error) error {
	err := run()
	if err != nil {
		return err
	}
	err = errors.Join(first, second)
	return err
}

func reuseAfterConditionalHandling(run func() error, terminated bool, collected []error) error {
	err := run()
	if err != nil && !terminated {
		return err
	}
	err = errors.Join(collected...)
	return err
}

func discardUnhandledError(run func() error, collected []error) error {
	err := run()
	err = errors.Join(collected...) // want "original error is discarded"
	return err
}

func discardErrorAfterReassignment(run func() error, collected []error) error {
	err := run()
	if err != nil {
		return err
	}
	err = run()
	err = errors.Join(collected...) // want "original error is discarded"
	return err
}

func unrelatedGuard(run func() error, collected []error) error {
	err := run()
	if err := run(); err != nil {
		return err
	}
	err = errors.Join(collected...) // want "original error is discarded"
	return err
}

func preserveParenthesizedError(run func() error, other error) error {
	err := run()
	err = errors.Join((err), other)
	return err
}
