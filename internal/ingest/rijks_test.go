package ingest

import (
	"encoding/xml"
	"fmt"
	"github.com/matryer/is"
	"net/http"
	"net/http/httptest"
	"rijks/internal/models"
	"testing"
)

var getRecordResponseSuccessful = &models.GetRecordResponse{
	ResponseDate: "2024-08-23T15:18:02Z",
	Request: models.Request{
		Verb:           "GetRecord",
		Identifier:     "oai:rijksmuseum.nl:sk-c-5",
		MetadataPrefix: "dc",
		Value:          "https://www.rijksmuseum.nl/api/oai/APIKEY",
	},
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
					Description: "Rembrandts beroemdste en grootste doek werd gemaakt voor de Kloveniersdoelen. Dit was een van de verenigingsgebouwen van de Amsterdamse schutterij, de burgerwacht van de stad. Rembrandt was de eerste die op een groepsportret de figuren in actie weergaf. De kapitein, in het zwart, geeft zijn luitenant opdracht dat de compagnie moet gaan marcheren. De schutters stellen zich op. Met behulp van licht vestigde Rembrandt de aandacht op belangrijke details, zoals het handgebaar van de kapitein en het kleine meisje op de achtergrond. Zij is de mascotte van de schutters.",
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

func TestRijksHandlder_GetRecord(t *testing.T) {
	// this test will make a call to a mock API using http test pacakge
	// then it will assert that we get back the intended response from that API.
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

	rh := NewRijksHandler("apikey", server.URL, client)
	got, err := rh.GetRecord("someID")
	is.NoErr(err)
	is.Equal(got.GetRecord.Record, getRecordResponseSuccessful.GetRecord.Record)
}
