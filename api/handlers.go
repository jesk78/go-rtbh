package api

import (
	"encoding/json"
	"net/http"
)

type greetings struct {
	Intro    string
	Messages []string
}

type blacklistApiResponse struct {
	Address string `json:"address"`
	Reason  string `json:"reason"`
}

type blacklistResponse struct {
	Entries []blacklistApiResponse `json:"entries"`
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	passedObj := greetings{
		Intro:    "Hello from Go!",
		Messages: []string{"Hello!", "Hi!", "¡Hola!", "Bonjour!", "Ciao!", "<script>evilScript()</script>"},
	}
	templates.ExecuteTemplate(w, "homePage", passedObj)
}

func blacklistHandler(w http.ResponseWriter, r *http.Request) {
	response := blacklistResponse{}

	for _, address := range Blacklist.GetAll() {
		entry := blacklistApiResponse{
			Address: address,
			Reason:  Blacklist.Get(address),
		}
		response.Entries = append(response.Entries, entry)
	}
	templates.ExecuteTemplate(w, "blacklist", response)
}

/*
 * Handler for GET /v1/blacklist
 */
func getBlacklist(w http.ResponseWriter, r *http.Request) {
	var entries []blacklistApiResponse
	var data []byte
	var err error

	for _, address := range Blacklist.GetAll() {
		entry := blacklistApiResponse{
			Address: address,
			Reason:  Blacklist.Get(address),
		}
		entries = append(entries, entry)
	}

	if data, err = json.Marshal(entries); err != nil {
		Log.Warning("[api.getBlacklist]: json.Marshal failed: " + err.Error())
		return
	}

	w.Write(data)
}
