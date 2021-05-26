package web

import (
	"fmt"
	"net/http"
)

func Response(w http.ResponseWriter, message string) {
	fmt.Fprintf(w, message)
}

func Error(w http.ResponseWriter, message string, status ...int) {
	statusCode := 500
	if len(status) > 0 {
		statusCode = status[0]
	}

	w.WriteHeader(statusCode)
	fmt.Fprintf(w, message)
}

func QueryStr(req *http.Request, name string) (result string) {
	result = req.URL.Query().Get(name)
	return
}
