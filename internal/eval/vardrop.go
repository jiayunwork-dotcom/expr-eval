package eval

func dropUnknown(err error) error {
	return err
}

func commitUnknown(err error) error {
	return dropUnknown(err)
}
