package es

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func (e *ES) FetchDetails(query io.Reader) ([]byte, error) {
	uri := fmt.Sprintf("http://%s/_search", e.cfg.ES.Address)
	request, err := http.NewRequest("POSt", uri, query)
	if err != nil {
		return nil, fmt.Errorf("ES.FetchDetails http.NewRequest: %v", err)
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accepts", "application/json")

	response, err := e.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("ES.FetchDetails client.Do: %v", err)
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("ES.FetchDetails ioutil.ReadAll: %v", err)
	}

	return data, nil
}
