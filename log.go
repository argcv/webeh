package webeh

import (
	"fmt"
	"net/http"
	"github.com/argcv/webeh/log"
)

func Log(v ...interface{}) {
	log.Infod(1, v...)
}

func Logf(f string, v ...interface{}) {
	log.Infod(1, fmt.Sprintf(f, v...))
}

func WebLog(r *http.Request, v ...interface{}) {
	log.Infod(2, fmt.Sprintf("%s - \"%s %s\" - \"%v\" - %v",
		GetUserIp(r), r.Method, r.URL,
		fmt.Sprint(v...), r.UserAgent()))
}

func WebLogf(r *http.Request, f string, v ...interface{}) {
	log.Infod(2, fmt.Sprintf("%s - \"%s %s\" - \"%v\" - %v",
		GetUserIp(r), r.Method, r.URL,
		fmt.Sprintf(f, v...), r.UserAgent()))
}
