package birdapi

import (
	"encoding/json"
	"fmt"
	"github.com/r3boot/go-rtbh/proto"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const TF_CLF string = "02/Jan/2006:15:04:05 -0700"

func HttpLog(r *http.Request, code int, size int) {
	var srcip string
	var logline string

	srcip = r.Header.Get("X-Forwarded-For")
	if srcip == "" {
		srcip = r.RemoteAddr
	}

	logline = srcip + " - - [" + time.Now().Format(TF_CLF) + "] "
	logline = logline + "\"" + r.Method + " " + r.URL.Path + " " + r.Proto + "\" "
	logline = logline + strconv.Itoa(code) + " " + strconv.Itoa(size)
	fmt.Println(logline)
}

// Default http response
func HttpResponse(w http.ResponseWriter, r *http.Request, msg string) {
	var data []byte
	var err error

	if data, err = json.Marshal(msg); err != nil {
		msg = "[HttpResponse] Failed to marshal message: " + err.Error()
	}

	w.Write(data)
	HttpLog(r, 200, len(data))
}

func HttpResponseData(w http.ResponseWriter, r *http.Request, data []byte) {
	w.Write(data)
	HttpLog(r, 200, len(data))
}

// <default http handler>
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	var msg string

	msg = "Nothing to see here, move along"

	HttpResponse(w, r, msg)
}

// GET /v1/status -- Gets an overview of all configured protocols
func getStatus(w http.ResponseWriter, r *http.Request) {
	var msg string
	var data []byte
	var output []string
	var fields []string
	var header string
	var status proto.BirdStatusOutput
	var statusline proto.BirdStatusProtocolEntry
	var err error

	if output, err = Bird.Command("show protocols"); err != nil {
		msg = "Failed to retrieve protocol output: " + err.Error()
		HttpResponse(w, r, msg)
		return
	}

	status = proto.BirdStatusOutput{}
	for _, line := range output {
		if len(line) > 5 {
			header = line[0:4]
			if header == "BIRD" {
				continue
			} else if header == "name" {
				continue
			}

			fields = strings.Fields(line)
			statusline = proto.BirdStatusProtocolEntry{
				Name:     fields[0],
				Protocol: fields[1],
				Table:    fields[2],
				State:    fields[3],
				Since:    fields[4],
				Info:     strings.Join(fields[5:], " "),
			}
			status.Protocols = append(status.Protocols, statusline)
		}
	}

	if data, err = json.Marshal(status); err != nil {
		msg = "Failed to marshal json output: " + err.Error()
		HttpResponse(w, r, msg)
		return
	}

	HttpResponseData(w, r, data)
}

func updateConfig(w http.ResponseWriter, r *http.Request) {
	var msg string
	var blacklist []string
	var whitelist []string

	blacklist = Blacklist.GetAll()
	whitelist = Whitelist.GetAll()

	if Bird.ExportPrefixes(whitelist, blacklist) {
		msg = "ok"
	} else {
		msg = "failed"
	}

	HttpResponse(w, r, msg)
}
