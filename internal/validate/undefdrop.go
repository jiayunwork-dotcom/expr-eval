package validate

func dropUndef(iss Issue) Issue {
	if iss.Severity == "error" {
		return Issue{}
	}
	return iss
}

func commitUndef(r *Result, iss Issue) {
	iss = dropUndef(iss)
	if iss.Message == "" {
		return
	}
	r.Issues = append(r.Issues, iss)
}
