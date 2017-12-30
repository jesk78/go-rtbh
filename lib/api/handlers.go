package api

import (
	"encoding/json"
	"github.com/r3boot/go-rtbh/lib/events"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type apiBlacklistEntries struct {
	data []*events.APIEvent `json:"data"`
}

/*
 * Various handlers for in-app paths
 */
func (api *RtbhApi) redirectToHomepage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", 301)
}

/* GET /app/*
 */
func (api *RtbhApi) handleAppDir() http.Handler {
	var app_dir http.Dir

	app_dir = http.Dir(Config.Api.Resources + "/app/")
	return http.StripPrefix("/app/", http.FileServer(app_dir))
}

/* GET /node_modules/*
 */
func (api *RtbhApi) handleNodeModulesDir() http.Handler {
	var node_dir http.Dir

	node_dir = http.Dir(Config.Api.Resources + "/node_modules/")
	return http.StripPrefix("/node_modules/", http.FileServer(node_dir))
}

/* GET /
 */
func (api *RtbhApi) handleHomepage() http.Handler {
	var homepage http.Dir

	homepage = http.Dir(Config.Api.Resources + "/index.html")
	return http.FileServer(homepage)
}

/*
 * Handler for various static assets
 */
func (api *RtbhApi) handleFileRequest(w http.ResponseWriter, r *http.Request) {
	var fs os.FileInfo
	var fd *os.File
	var file_size int64
	var bytes_read int
	var data []byte
	var full_path string
	var err error

	if r.URL.Path == "/" {
		full_path = Config.Api.Resources + "/index.html"
	} else {
		full_path = Config.Api.Resources + r.URL.Path
	}

	if fs, err = os.Stat(full_path); err != nil {
		w.WriteHeader(http.StatusNotFound)
		Log.Warning(MYNAME + ": " + err.Error())
		return
	}

	if fs.IsDir() {
		w.WriteHeader(http.StatusNotFound)
		Log.Warning(MYNAME + ": " + full_path + ": is a directory")
		return
	}

	file_size = fs.Size()

	if fd, err = os.Open(full_path); err != nil {
		w.WriteHeader(http.StatusNotFound)
		Log.Warning(MYNAME + ": Failed to open " + full_path + ": " + err.Error())
		return
	}

	data = make([]byte, file_size)
	if bytes_read, err = fd.Read(data); err != nil {
		w.WriteHeader(http.StatusNotFound)
		Log.Warning(MYNAME + ": Failed to read " + full_path + ": " + err.Error())
		return
	}

	if int64(bytes_read) != file_size {
		w.WriteHeader(http.StatusNotFound)
		Log.Warning(MYNAME + ": Only read " + strconv.Itoa(bytes_read) + " out of " + strconv.Itoa(int(file_size)) + " bytes")
		return
	}

	w.Write(data)
}

/*
 * Handler for /v1/blacklist[:id]
 */
func (api *RtbhApi) handleBlacklist(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			api.getBlacklist(w, r)
		}
	case http.MethodPut:
		{
			Log.Debug(MYNAME + ": Handling save request")
			api.saveBlacklist(w, r)
		}
	default:
		{
			w.WriteHeader(http.StatusNotFound)
			Log.Warning(MYNAME + ": Unknown request method used for " + PATH_API_BLACKLIST)
		}
	}
}

/*
 * Handler for GET /v1/blacklist[:<id>]
 */
func (api *RtbhApi) getBlacklist(w http.ResponseWriter, r *http.Request) {
	var id int64
	var id_s string
	var entries []*events.APIEvent
	var entry *events.APIEvent
	var data []byte
	var err error

	if strings.Contains(r.URL.Path, ":") {
		id_s = strings.Split(r.URL.EscapedPath(), ":")[1]
		if id, err = strconv.ParseInt(id_s, 10, 64); err != nil {
			w.WriteHeader(http.StatusNotFound)
			Log.Warning(MYNAME + ".getBlacklist: " + err.Error())
			return
		}

		if entry = api.blacklist.GetById(id); entry == nil {
			w.WriteHeader(http.StatusNotFound)
			Log.Warning(MYNAME + ".getBlacklist: Entry for " + strconv.Itoa(int(id)) + " not found")
			return
		}

		entries = append(entries, entry)
	} else {
		entries = api.blacklist.GetAll()
	}

	if data, err = json.Marshal(entries); err != nil {
		w.WriteHeader(http.StatusNotFound)
		Log.Warning(MYNAME + ".getBlacklist: json.Marshal failed: " + err.Error())
		return
	}

	w.Write(data)
}

/*
 * Handler for PUT /v1/blacklist[:<id>]
 */
func (api *RtbhApi) saveBlacklist(w http.ResponseWriter, r *http.Request) {
	var id_s string
	var id int64
	var err error

	if strings.Contains(r.URL.Path, ":") {
		id_s = strings.Split(r.URL.EscapedPath(), ":")[1]
		if id, err = strconv.ParseInt(id_s, 10, 64); err != nil {
			w.WriteHeader(http.StatusNotFound)
			Log.Warning(MYNAME + ".getBlacklist: " + err.Error())
			return
		}

	}

	Log.Debug(MYNAME + ".saveBlacklist: Saving id " + strconv.Itoa(int(id)))
}
