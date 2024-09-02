package ingest

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"rijks/internal/models"
)

const (
	getRecordURL   = "?verb=GetRecord&metadataPrefix=dc&identifier="
	listRecordsURL = "?verb=ListRecords&set=subject:EntirePublicDomainSet&metadataPrefix=dc"
)

var ErrRecordNotFound = errors.New("record not found")

type RijksHandler struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewRijksHandler(apiKey, baseURL string, client *http.Client) *RijksHandler {
	return &RijksHandler{apiKey: apiKey, baseURL: baseURL, client: client}
}

func (rh *RijksHandler) SetAPIKey(apiKey string) {
	rh.apiKey = apiKey
}

func (rh *RijksHandler) buildGetRecordURL(identifier string) (*url.URL, error) {
	requestURL, err := url.Parse(fmt.Sprintf("%s/%s%s%s", rh.baseURL, rh.apiKey, getRecordURL, identifier))
	if err != nil {
		return nil, err
	}
	return requestURL, nil
}

func (rh *RijksHandler) buildListRecordsURL() (*url.URL, error) {
	requestURL, err := url.Parse(fmt.Sprintf("%s/%s%s", rh.baseURL, rh.apiKey, listRecordsURL))
	if err != nil {
		return nil, err
	}
	return requestURL, nil
}

func (rh *RijksHandler) GetRecord(identifier string) (*models.GetRecordResponse, error) {
	// Build the requestURL.
	requestURL, err := rh.buildGetRecordURL(identifier)
	if err != nil {
		return nil, err
	}

	// Construct the request object.
	req, err := http.NewRequest(http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return nil, err
	}

	// Handle the request.
	response, err := rh.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status not ok %d", response.StatusCode)
	}
	//Read response body.
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshall the response body.
	getRecordResponse := &models.GetRecordResponse{}
	err = xml.Unmarshal(body, getRecordResponse)
	if err != nil {
		return nil, err
	}
	if getRecordResponse.Error != nil {
		if getRecordResponse.Error.Code == "idDoesNotExist" {
			return nil, fmt.Errorf("GetRecord: %s %w", identifier, ErrRecordNotFound)
		}
		return nil, fmt.Errorf("GetRecord: %s %v", identifier, getRecordResponse.Error)
	}

	return getRecordResponse, nil
}

func (rh *RijksHandler) ListRecords() (*models.ListRecordsResponse, error) {
	requestURL, err := rh.buildListRecordsURL()
	if err != nil {
		return nil, err
	}

	// Construct the request object.
	req, err := http.NewRequest(http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return nil, err
	}

	// Handle the request.
	response, err := rh.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status not ok %d", response.StatusCode)
	}

	//Read response body.
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(body)
	lrr := &models.ListRecordsResponse{}
	err = xml.Unmarshal(body, lrr)

	return nil, err

}
