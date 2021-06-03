package web

import (
	"bank-api/pkg"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"
)

func Response(w http.ResponseWriter, message string) {
	fmt.Fprintf(w, "%v", message)
	// fmt.Fprintf(w, message)
}

func Error(w http.ResponseWriter, message interface{}, status ...int) {
	statusCode := 500
	if len(status) > 0 {
		statusCode = status[0]
	}
	w.WriteHeader(statusCode)

	s := message.(string)
	fmt.Fprintf(w, "%v", s)
	// fmt.Fprintf(w, s)

}

func QueryStr(req *http.Request, name string) (result string) {
	result = req.URL.Query().Get(name)
	return
}
func QueryStrToDecimal(r *http.Request, name string) decimal.Decimal {
	return pkg.StringToDecimal(r.URL.Query().Get(name))
}
