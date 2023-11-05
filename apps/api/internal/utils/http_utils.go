package utils

import "net/http"

func WriteText(writer http.ResponseWriter, str string) {
	_, err := writer.Write([]byte(str))
	if err != nil {
		writer.WriteHeader(500)
		return
	}
}
