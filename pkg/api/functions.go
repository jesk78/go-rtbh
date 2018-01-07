package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type logWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (r WebResponse) ToJSON() []byte {
	data, err := json.Marshal(r)
	if err != nil {
		return []byte("{\"status\":false,\"message\":\"json encoding failed\"}")
	}

	return data
}

func (w *logWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}

	w.length = len(b)

	return w.ResponseWriter.Write(b)
}

func (w *logWriter) WriteHeader(status int) {
	w.ResponseWriter.Header().Add("Content-Length", strconv.Itoa(w.length))
	w.ResponseWriter.WriteHeader(status)
	w.status = status
}

func (a *RtbhApi) SetupRouting() {
	// API
	http.HandleFunc("/api/v1/blacklist", a.BlacklistHandler)
	http.HandleFunc("/api/v1/whitelist", a.WhitelistHandler)

	// Website
	http.Handle("/css/", logHandler(http.FileServer(http.Dir("./ui"))))
	http.Handle("/js/", logHandler(http.FileServer(http.Dir("./ui"))))
	http.Handle("/img/", logHandler(http.FileServer(http.Dir("./ui"))))
	http.Handle("/fonts/", logHandler(http.FileServer(http.Dir("./ui"))))
	http.HandleFunc("/", a.HomeHandler)
}

func (api *RtbhApi) Run() {
	var url string

	url = cfg.Api.BindIp + ":" + cfg.Api.BindPort
	log.Debugf("RtbhApi.Run: Listening on http://" + url)
	http.ListenAndServe(url, nil)

	return
}

func (api *RtbhApi) HttpLog(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var custom_writer *logWriter
		var srcip string
		var logline string
		var t_start time.Time
		var latency_d time.Duration
		var latency float64
		var tokens []string

		t_start = time.Now()
		custom_writer = &logWriter{w, 0, 0}
		handler.ServeHTTP(custom_writer, r)
		latency_d = time.Now().Sub(t_start)
		latency = float64(latency_d.Nanoseconds()) / 1000000

		// Get (proxied) source ip address
		srcip = r.Header.Get("X-Forwarded-For")
		if srcip == "" {
			if r.RemoteAddr == "" {
				srcip = "unknown"
			} else {
				srcip = r.RemoteAddr
			}
		}

		// Split port if it is appended
		if strings.Contains(srcip, ":") {
			tokens = strings.Split(srcip, ":")
			if len(tokens) == 2 {
				// IPv4 address
				srcip = strings.Join(tokens[:len(tokens)-1], ":")
			} else {
				srcip = strings.Join(tokens[:len(tokens)-1], ":")[1:]
				srcip = srcip[:len(srcip)-1]
			}
		}

		logline = srcip + " - - [" + time.Now().Format(TF_CLF) + "] "
		logline = logline + "\"" + r.Method + " " + r.URL.Path + " " + r.Proto + "\" "
		logline = logline + strconv.Itoa(custom_writer.status) + " "
		logline = logline + strconv.Itoa(custom_writer.length) + " "
		logline = logline + fmt.Sprintf("%.02f", latency)

		fmt.Println(logline)
	})
}
