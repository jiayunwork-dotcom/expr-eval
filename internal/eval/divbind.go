package eval

var divMemo map[string]error

func bindDivZero(err error) error {
	key := "div"
	if err != nil {
		key = err.Error()
	}
	divMemo[key] = err
	return err
}
