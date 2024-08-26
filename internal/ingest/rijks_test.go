package ingest

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"

	"rijks/internal/models"
)

var (
	apiKey                    = "apiKey"
	identifier                = "someID"
	getRecordResponseNotFound = &models.GetRecordResponse{Error: &models.Error{
		Code:    "idDoesNotExist",
		Message: "Record Not Found",
	},
	}

	getRecordResponseUnexpectedError = &models.GetRecordResponse{Error: &models.Error{
		Code:    "someCode",
		Message: "NO!",
	},
	}
	getRecordResponseSuccessful = &models.GetRecordResponse{

		GetRecord: &models.GetRecord{
			Record: models.Record{
				Header: models.Header{
					Identifier: "oai:rijksmuseum.nl:SK-C-5",
					Datestamp:  "2024-07-23T15:37:03Z",
				},
				Metadata: models.Metadata{
					OaiDc: models.OaiDc{
						Identifiers: []string{
							"http://hdl.handle.net/10934/RM0001.COLLECT.5216",
							"SK-C-5",
						},
						Title:   "De Nachtwacht",
						Creator: "Rijn, Rembrandt van",
						Subjects: []string{
							"Amsterdam",
							"Banninck Cocq, Frans",
							"Ruytenburch, Willem van",
							"Visscher Cornelisen, Jan",
							"Kemp, Rombout",
							"Engelen, Reijnier Janszn",
							"Bolhamer, Barent Harmansen",
							"Keijser, Jan Adriaensen",
							"Willemsen, Elbert",
							"Leijdeckers, Jan Claesen",
							"Ockersen, Jan",
							"Bronchorst, Jan Pietersen",
							"Wormskerck, Harman Jacobsen",
							"Roy, Jacob Dircksen de",
							"Heede, Jan van der",
						},
						Description: "Rembrandts beroemdste en grootste doek werd gemaakt voor de Kloveniersdoelen.",
						Date:        "1642",
						Type:        "schilderij",
						Formats:     []string{"doek", "olieverf"},
						Language:    "nl",
						Publisher:   "Rijksmuseum",
						Rights:      "http://creativecommons.org/publicdomain/mark/1.0/",
						Coverage:    "Amsterdam",
					},
				},
			},
		},
	}
)

func TestRijksHandlder_GetRecord(t *testing.T) {
	// this test will make a call to a mock API using http test pacakge
	// then it will assert that we get back the intended response from that API.
	is := is.New(t)
	t.Run("returns the correct record when request is successful", func(t *testing.T) {
		is := is.New(t)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bytes, err := xml.Marshal(getRecordResponseSuccessful)
			if err != nil {
				t.Fatalf("unable to marshal xml: %v", err)
			}
			w.WriteHeader(http.StatusOK)
			_, err = fmt.Fprintln(w, string(bytes))
			if err != nil {
				t.Fatalf("unable to write response: %v", err)
			}

		}))
		defer server.Close()
		client := &http.Client{}

		rh := NewRijksHandler(apiKey, server.URL, client)
		got, err := rh.GetRecord(identifier)
		is.NoErr(err)
		is.Equal(got.GetRecord.Record, getRecordResponseSuccessful.GetRecord.Record)
	})

	t.Run("returns an error when the record is not found", func(t *testing.T) {
		is := is.New(t)

		expectedErr := fmt.Errorf("GetRecord: %s %w", identifier, ErrRecordNotFound)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bytes, err := xml.Marshal(getRecordResponseNotFound)
			if err != nil {
				t.Fatalf("unable to marshal struct")
			}

			w.WriteHeader(http.StatusOK)
			_, err = fmt.Fprintln(w, string(bytes))
			if err != nil {
				t.Fatalf("unable to write response: %v", err)
			}
		}))

		defer server.Close()
		client := &http.Client{}

		rh := NewRijksHandler(apiKey, server.URL, client)
		got, err := rh.GetRecord(identifier)
		is.Equal(err, expectedErr)
		is.Equal(got, nil)
	})

	t.Run("returns an error when response contains an unexpected error", func(t *testing.T) {
		is := is.New(t)

		expectedErr := fmt.Errorf("GetRecord: %s %v", identifier, getRecordResponseUnexpectedError.Error)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bytes, err := xml.Marshal(getRecordResponseUnexpectedError)
			if err != nil {
				t.Fatalf("unable to marshal struct")
			}

			w.WriteHeader(http.StatusOK)
			_, err = fmt.Fprintln(w, string(bytes))
			if err != nil {
				t.Fatalf("unable to write response: %v", err)
			}
		}))

		defer server.Close()
		client := &http.Client{}

		rh := NewRijksHandler(apiKey, server.URL, client)
		got, err := rh.GetRecord(identifier)
		is.Equal(err, expectedErr)
		is.Equal(got, nil)
	})

}
