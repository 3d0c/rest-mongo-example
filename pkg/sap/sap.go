package sap

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/teal-seagull/lyre-be-v4/pkg/config"
)

const (
	// @TODO move to config
	timeout = time.Second * 15
)

// Request structure
type Request struct {
	*http.Request
	ctx    context.Context
	cancel context.CancelFunc
}

// NewRequest returns initialized SAP request
func NewRequest(method, url string, body []byte) (*Request, error) {
	var (
		req *Request = &Request{
			Request: new(http.Request),
		}
		err error
	)

	req.ctx, req.cancel = context.WithTimeout(context.Background(), timeout)

	if len(body) == 0 {
		req.Request, err = http.NewRequestWithContext(req.ctx, method, url, nil)
	} else {
		req.Request, err = http.NewRequestWithContext(req.ctx, method, url, bytes.NewReader(body))
	}
	if err != nil {
		return nil, fmt.Errorf("error creating request - %s", err)
	}

	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Authorization", config.TheConfig().SAP.Auth)

	return req, nil
}

// Do implements http.Do request
func (r *Request) Do(result interface{}) error {
	var (
		resp   *http.Response
		client http.Client
		err    error
	)
	defer r.cancel()

	if resp, err = client.Do(r.Request); err != nil {
		return fmt.Errorf("error doing request - %s", err)
	}

	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("error decoding response - %s", err)
	}

	return nil
}

// ParseTimeStamp parses SAP timestamp
func ParseTimeStamp(s string) (time.Time, error) {
	return time.Parse("20060102150405", s)
}
