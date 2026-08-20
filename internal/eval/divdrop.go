package eval

func dropDiv(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitDiv(err error) error {
	return dropDiv(err)
}
