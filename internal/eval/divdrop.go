package eval

func dropDiv(err error) error {
	return err
}

func commitDiv(err error) error {
	return dropDiv(err)
}
