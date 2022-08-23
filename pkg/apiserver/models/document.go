package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/teal-seagull/lyre-be-v4/pkg/config"
	"github.com/teal-seagull/lyre-be-v4/pkg/sap"
)

// DocumentScheme structure
type DocumentScheme struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	DocNumber   string    `json:"doc_number"`
	Version     string    `json:"version"`
	Part        string    `json:"part"`
	CreatedBy   string    `json:"created_by"`
	UpdatedBy   string    `json:"updated_by"`
	CreatedDate time.Time `json:"created_date,omitempty"`
	UpdatedDate time.Time `json:"updated_date,omitempty"`
}

// Document model
type Document struct{}

// NewDocument document model constructor
func NewDocument() (*Document, error) {
	return &Document{}, nil
}

// Find documents
func (d *Document) Find(item string, typ string) ([]DocumentScheme, error) {
	var (
		req  *sap.Request
		list SAPDocList
		err  error
	)

	if req, err = sap.NewRequest(
		"GET",
		config.TheConfig().SAP.DocList,
		nil,
	); err != nil {
		return nil, fmt.Errorf("error creating SAP request - %s", err)
	}

	q := req.URL.Query()
	q.Add("ip_docid", item)
	q.Add("ip_dokob", typ)

	req.URL.RawQuery = q.Encode()

	if err = req.Do(&list); err != nil {
		return nil, fmt.Errorf("error doing request - %s", err)
	}

	return list.ToDocumentList(), nil
}

// Download creates a request to SAP to download file by ID to the local server to provided path
// It returns requested file name if it has been downloaded or empty if not found
func (d *Document) Download(id string, path string) (string, error) {
	type request struct {
		ID   string `json:"ip_fileid"`
		Path string `json:"ip_serverpath"`
	}
	type access struct {
		FileName   string `json:"COMP_ID"`
		FileSize   string `json:"COMP_SIZE"`
		MimeType   string `json:"MIMETYPE"`
		CreatedAt  string `json:"CREA_TIME"`
		ChangedAt  string `json:"CHNG_TIME"`
		BinaryFlag string `json:"BINARY_FLG"`
		FirstLine  int    `json:"FIRST_LINE"`
		LastLine   int    `json:"LAST_LINE"`
	}
	type SAPresult struct {
		Access []access `json:"ET_ACCESS"`
	}

	var (
		req     *sap.Request
		list    SAPDocList
		result  SAPresult
		payload []byte
		err     error
	)

	if payload, err = json.Marshal(request{ID: id, Path: path}); err != nil {
		return "", fmt.Errorf("error creating request payload - %s", err)
	}

	if req, err = sap.NewRequest(
		"PUT",
		config.TheConfig().SAP.DocGet,
		payload,
	); err != nil {
		return "", fmt.Errorf("error creating SAP request - %s", err)
	}

	if err = req.Do(&list); err != nil {
		return "", fmt.Errorf("error doing request - %s", err)
	}

	if len(result.Access) == 0 {
		// @TODO Custom error here
		return "", nil
	}

	return result.Access[0].FileName, nil
}

// SAPDoc is a structure for mapping SAP response
type SAPDoc struct {
	DocumentType  string `json:"DOKAR"`
	DocumentID    string `json:"DOKNR"`
	DocumentVer   string `json:"DOKVR"`
	DocumentPart  string `json:"DOKTL"`
	ObjectLink    string `json:"DOKOB"`
	FileID        string `json:"FILEID"`
	WsApplication string `json:"WSAPPLICATION"`
	Mimetype      string `json:"MIMETYPE"`
	DocFile       string `json:"DOCFILE"`
	CreatedAt     string `json:"CREATED_AT"`
	ChangedAt     string `json:"CHANGED_AT"`
	CreatedBy     string `json:"CREATED_BY"`
	ChangedBy     string `json:"CHANGED_BY"`
}

// ToDocument maps SAPDoc into Document
func (doc SAPDoc) ToDocument() DocumentScheme {
	var (
		result = DocumentScheme{}
	)

	result.ID = doc.FileID
	result.DocNumber = doc.DocumentID
	result.Name = doc.DocFile
	result.Type = doc.WsApplication
	result.Version = doc.DocumentVer
	result.Part = doc.DocumentPart
	result.CreatedBy = doc.CreatedBy
	result.UpdatedBy = doc.ChangedBy

	// Ingoring error, CreatedDate will be just empty
	result.CreatedDate, _ = sap.ParseTimeStamp(doc.CreatedAt)
	result.UpdatedDate, _ = sap.ParseTimeStamp(doc.ChangedAt)

	return result
}

// SAPDocList is a structure for mapping SAP response
type SAPDocList struct {
	Output []SAPDoc `json:"ET_OUTPUT"`
}

// ToDocumentList maps list of SAP response into array of DocumentSchemes
func (list SAPDocList) ToDocumentList() []DocumentScheme {
	var (
		result = make([]DocumentScheme, 0)
	)

	for i := range list.Output {
		result = append(result, list.Output[i].ToDocument())
	}

	return result
}
