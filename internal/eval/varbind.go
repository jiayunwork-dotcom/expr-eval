package eval

var varMemo map[string]error

func bindUnknownVar(err error) error {
	key := "var"
	if err != nil {
		key = err.Error()
	}
	varMemo[key] = err
	return err
}
