package testdata

func _() error {
	var err error
	if err != nil {
		return nil // want ".+"
	}

	if err != nil {
		return err
	}

	if err == nil {
		return nil
	}

	return nil
}

func _() (int, error) {
	var err error
	if err != nil {
		return 0, nil // want ".+"
	}

	if err != nil {
		return 0, err
	}

	if err == nil {
		return 0, nil
	}

	return 0, nil
}

func _() (*int, error) {
	var err error
	if err != nil {
		return nil, nil // want ".+"
	}

	if err != nil {
		return nil, err
	}

	if err == nil {
		return nil, nil
	}

	return nil, nil
}

func _() {
	var err error
	if err != nil {
		return
	}
	return
}

func _() *int {
	var err error
	if err != nil {
		return nil
	}
	return nil
}
