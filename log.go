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

func WebLog(r *http.Request, v ...interface{}) {
	logger.Output(2, fmt.Sprintf("%s - \"%s %s\" - \"%v\" - %v",
		GetUserIp(r), r.Method, r.URL,
		fmt.Sprint(v...), r.UserAgent()))
}

func init() {
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
}
