package webeh

import (
	"net"
	"net/http"
)

func GetUserIp(r *http.Request) string {
	agentRemoteHeaderUpyunKey := "Client-IP"
	agentRemoteHeaderKey := "X-Real-IP"
	if ip := r.Header.Get(agentRemoteHeaderUpyunKey); ip != "" {
		return ip
	} else if ip := r.Header.Get(agentRemoteHeaderKey); ip != "" {
		return ip
	} else {
		if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			return ip
		} else {
			return "0.0.0.0"
		}
	}
}
