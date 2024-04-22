package helper

import (
	"fmt"
	"net/http"

	"github.com/codecrafters-io/http-server-starter-go/pkg/constants"
)

func NewResponse(status int, body []byte, ct string) []byte {
	res := []byte{}
	statusText := "OK"
	switch status {
	case http.StatusNotFound:
		statusText = "Not Found"
	case http.StatusInternalServerError:
		statusText = "Internal Server Error"
	}

	res = append(res, []byte(fmt.Sprintf("HTTP/1.1 %d %v%s", status, statusText, constants.CRLF))...)

	if status == http.StatusNotFound || status == http.StatusInternalServerError {
		return append(res, constants.CRLF...)
	}

	if ct != "" {
		res = append(res, []byte(fmt.Sprintf("Content-Type: %s%s", ct, constants.CRLF))...)
	}

	size := len(body)
	if size != 0 {
		res = append(res, []byte(fmt.Sprintf("Content-Length: %d%s%s", size, constants.CRLF, constants.CRLF))...)
		res = append(res, []byte(fmt.Sprintf("%s%s", body, constants.CRLF))...)
	}

	return res
}
