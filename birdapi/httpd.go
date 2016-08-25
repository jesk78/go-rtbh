package birdapi

import (
	"net/http"
)

func ConfigureRouting() {
	Log.Debug("[routing] '/v1/update' -> updateConfig")
	http.HandleFunc("/v1/update", updateConfig)

	Log.Debug("[routing] '/v1/status' -> getStatus")
	http.HandleFunc("/v1/status", getStatus)

	Log.Debug("[routing] '/' -> defaultHandler")
	http.HandleFunc("/", defaultHandler)
}

func RunServer() {
	var addr string

	addr = Config.Api.BindIp + ":" + Config.Api.BindPort

	Log.Debug("[go-birdapi]: Listening on " + addr)
	http.ListenAndServe(addr, nil)
}
