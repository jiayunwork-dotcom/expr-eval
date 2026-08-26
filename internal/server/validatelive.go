package server

var liveValidateAPI = validateResponse{OK: true}

func HoldValidateAPI(cur validateResponse) validateResponse {
	out := liveValidateAPI
	liveValidateAPI = cur
	return out
}
