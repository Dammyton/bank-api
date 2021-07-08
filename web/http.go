package web

import (
	"bank-api/pkg"
	"fmt"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"
)

func Response(w http.ResponseWriter, message interface{}) {
	fmt.Fprintf(w, "%v", message)
}

func Error(w http.ResponseWriter, message interface{}, status ...int) {
	statusCode := 500
	if len(status) > 0 {
		statusCode = status[0]
	}
	w.WriteHeader(statusCode)

	fmt.Fprintf(w, "%v", message)

}

func QueryStr(req *http.Request, name string) (result int64) {
	s := req.URL.Query().Get(name)
	result, _ = strconv.ParseInt(s, 10, 64)

	return
}
func QueryStrToDecimal(r *http.Request, name string) decimal.Decimal {
	return pkg.StringToDecimal(r.URL.Query().Get(name))
}
