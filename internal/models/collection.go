package models

import "encoding/xml"

// GetRecordResponse represents the root element for the
// getRecord api from the rijks API
type GetRecordResponse struct {
	XMLName      xml.Name   `xml:"OAI-PMH"`
	ResponseDate string     `xml:"responseDate"`
	Request      Request    `xml:"request"`
	GetRecord    *GetRecord `xml:"GetRecord,omitempty"`
	Error        *Error     `xml:"error,omitempty"`
}

type Error struct {
	Code    string `xml:"code,attr"`
	Message string `xml:",chardata"`
}

// Request represents the request element
type Request struct {
	Verb           string `xml:"verb,attr"`
	Identifier     string `xml:"identifier,attr"`
	MetadataPrefix string `xml:"metadataPrefix,attr"`
	Value          string `xml:",chardata"`
}

// GetRecord represents the GetRecord element
type GetRecord struct {
	Record Record `xml:"record"`
}

// Record represents the record element
type Record struct {
	Header   Header   `xml:"header"`
	Metadata Metadata `xml:"metadata"`
}

// Header represents the header element
type Header struct {
	Identifier string `xml:"identifier"`
	Datestamp  string `xml:"datestamp"`
}

// Metadata represents the metadata element
type Metadata struct {
	OaiDc OaiDc `xml:"http://www.openarchives.org/OAI/2.0/oai_dc/ dc"`
}

// OaiDc represents the oai_dc:dc element
type OaiDc struct {
	Identifiers []string `xml:"http://purl.org/dc/elements/1.1/ identifier"`
	Title       string   `xml:"http://purl.org/dc/elements/1.1/ title"`
	Creator     string   `xml:"http://purl.org/dc/elements/1.1/ creator"`
	Subjects    []string `xml:"http://purl.org/dc/elements/1.1/ subject"`
	Description string   `xml:"http://purl.org/dc/elements/1.1/ description"`
	Date        string   `xml:"http://purl.org/dc/elements/1.1/ date"`
	Type        string   `xml:"http://purl.org/dc/elements/1.1/ type"`
	Formats     []string `xml:"http://purl.org/dc/elements/1.1/ format"`
	Language    string   `xml:"http://purl.org/dc/elements/1.1/ language"`
	Publisher   string   `xml:"http://purl.org/dc/elements/1.1/ publisher"`
	Rights      string   `xml:"http://purl.org/dc/elements/1.1/ rights"`
	Coverage    string   `xml:"http://purl.org/dc/elements/1.1/ coverage"`
}

type ListRecordsResponse struct {
	Records         []Record `xml:"record"`
	ResumptionToken string   `xml:"resumptionToken"`
}
