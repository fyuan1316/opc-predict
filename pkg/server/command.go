package server

import "bytes"

var ControlCommand = struct {
	Start string
	Stop  string
}{
	"startIt",
	"stopIt",
}

var CommandPrefix = "CTRL:"

func IsPureData(data []byte) bool {
	if bytes.HasPrefix(data, []byte("CTRL:")) {
		return false
	}
	return true
}
