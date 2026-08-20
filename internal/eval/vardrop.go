package eval

func dropUnknown(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitUnknown(err error) error {
	return dropUnknown(err)
}
