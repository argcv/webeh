package webeh

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var logger *log.Logger

func Log(v ...interface{}) {
	logger.Output(2, fmt.Sprintln(v...))
}

func Logf(f string, v ...interface{}) {
	logger.Output(2, fmt.Sprintf(f, v...))
}

func WebLog(r *http.Request, v ...interface{}) {
	logger.Output(2, fmt.Sprintf("%s - \"%s %s\" - \"%v\" - %v",
		GetUserIp(r), r.Method, r.URL,
		fmt.Sprint(v...), r.UserAgent()))
}

func WebLogf(r *http.Request, f string, v ...interface{}) {
	logger.Output(2, fmt.Sprintf("%s - \"%s %s\" - \"%v\" - %v",
		GetUserIp(r), r.Method, r.URL,
		fmt.Sprintf(f, v...), r.UserAgent()))
}

func init() {
	logger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
}
