package api

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"

	"bytes"

	"github.com/r3boot/go-rtbh/pkg/events"
)

type apiBlacklistEntries struct {
	data []*events.APIEvent `json:"data"`
}

func logHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpLog(r, http.StatusOK, 0)
		h.ServeHTTP(w, r)
	})
}

func (a *RtbhApi) HomeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			content, err := ioutil.ReadFile("./ui/templates/index.html")
			if err != nil {
				a.log.Warningf("WebAPI.HomeHandler ioutil.ReadFile: %v", err)
				errorResponse(w, r, "Failed to read template")
				return
			}

			t := template.New("index")

			_, err = t.Parse(string(content))
			if err != nil {
				a.log.Warningf("WebAPI.HomeHandler t.Parse: %v", err)
				errorResponse(w, r, "Failed to parse template")
				return
			}

			data := TemplateData{}

			output := bytes.Buffer{}

			err = t.Execute(&output, data)
			if err != nil {
				a.log.Warningf("WebAPI.HomeHandler t.Execute: %v", err)
				errmsg := "Failed to execute template"
				http.Error(w, errmsg, http.StatusInternalServerError)
				httpLog(r, http.StatusInternalServerError, len(errmsg))
				return
			}

			w.Write(output.Bytes())
			httpLog(r, http.StatusOK, output.Len())
		}
	default:
		{
			errorResponse(w, r, "Unsupported method")
		}
	}
}

func (a *RtbhApi) BlacklistHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			entries, err := a.blacklist.GetAll()
			if err != nil {
				a.log.Warningf("RtbhApi.BlacklistHandler: %v", err)
				errorResponse(w, r, "Failed to retrieve blacklist entries")
				return
			}

			data, err := json.Marshal(&entries)
			if err != nil {
				a.log.Warningf("RtbhApi.BlacklistHandler json.Marshal: %v", err)
				errorResponse(w, r, "Failed to retrieve blacklist entries")
				return
			}

			w.Write(data)
			okResponse(r, len(data))
		}
	case http.MethodTrace:
		{

		}
	default:
		{
			errorResponse(w, r, "Unsupported method")
		}
	}
}

func (a *RtbhApi) WhitelistHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet: // Get all whitelist entries
		{
			entries, err := a.whitelist.GetAll()
			if err != nil {
				a.log.Warningf("RtbhApi.WhitelistHandler: %v", err)
				errorResponse(w, r, "Failed to retrieve whitelist entries")
				return
			}

			data, err := json.Marshal(&entries)
			if err != nil {
				a.log.Warningf("RtbhApi.Whitelist json.Marshal: %v", err)
				errorResponse(w, r, "Failed to retrieve whitelist entries")
				return
			}

			w.Write(data)
			okResponse(r, len(data))
		}
	case http.MethodPost: // Create new whitelist entry
		{
			decoder := json.NewDecoder(r.Body)

			request := &WebWhitelistAddRequest{}
			response := &WebResponse{}

			err := decoder.Decode(request)
			if err != nil {
				a.log.Warningf("RtbhApi.WhitelistHandler decoder.Decode: %v", err)
				errorResponse(w, r, "Failed to decode request")
				return
			}

			err = a.whitelist.Add(events.RTBHWhiteEntry{
				Address:     request.IpAddr,
				Description: request.Description,
			})
			if err != nil {
				a.log.Warningf("RtbhApi.WhitelistHandler: %v", err)
				errorResponse(w, r, "Failed to add entry")
				return
			}

			response.Status = true
			response.Message = "Added entry to whitelist"

			w.Write(response.ToJSON())
			okResponse(r, len(response.ToJSON()))
		}
	case http.MethodPatch:
		{
			decoder := json.NewDecoder(r.Body)

			request := &WebWhitelistAddRequest{}
			response := &WebResponse{}

			err := decoder.Decode(request)
			if err != nil {
				a.log.Warningf("RtbhApi.WhitelistHandler decoder.Decode: %v", err)
				errorResponse(w, r, "Failed to decode request")
				return
			}

			err = a.whitelist.Update(events.RTBHWhiteEntry{
				Address:     request.IpAddr,
				Description: request.Description,
			})
			if err != nil {
				a.log.Warningf("RtbhApi.WhitelistHandler: %v", err)
				errorResponse(w, r, "Failed to add entry")
				return
			}

			response.Status = true
			response.Message = "Updated whitelist entry"

			w.Write(response.ToJSON())
			okResponse(r, len(response.ToJSON()))
		}
	case http.MethodDelete:
		{
			decoder := json.NewDecoder(r.Body)

			request := &WebWhitelistRemoveRequest{}
			response := &WebResponse{}

			err := decoder.Decode(request)
			if err != nil {
				a.log.Warningf("RtbhApi.WhitelistHandler decoder.Decode: %v", err)
				errorResponse(w, r, "Failed to decode request")
				return
			}

			err = a.whitelist.Remove(request.IpAddr)
			if err != nil {
				a.log.Warningf("RtbhApi.WhitelistHandler: %v", err)
				errorResponse(w, r, "Failed to remove entry")
				return
			}

			response.Status = true
			response.Message = "Removed entry from whitelist"

			w.Write(response.ToJSON())
			okResponse(r, len(response.ToJSON()))
		}
	default:
		{
			errorResponse(w, r, "Unsupported method")
		}
	}
}

func (a *RtbhApi) ESProxyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		{
			details, err := a.es.FetchDetails(r.Body)
			if err != nil {
				a.log.Warningf("RtbhApi.ESProxyHandler: %v", err)
				errorResponse(w, r, "Failed to fetch details")
			}

			response := WebResponse{
				Status: true,
				Data:   details,
			}

			data := response.ToJSON()

			w.Write(data)
			okResponse(r, len(data))
		}
	default:
		{
			errorResponse(w, r, "Unsupported method")
		}
	}
}
