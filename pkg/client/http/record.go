package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ldassonville/happy-beer-api/pkg/api"
	"io"
	"net/http"
)

func (e *Client) SearchRecords(ctx context.Context) ([]*api.Record, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/records", e.config.ApiUrl), nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("an error happend while building request for listing records (cause: %s)", err.Error()))
	}

	req.Header.Set(headerNameAccept, headerValueApplicationJSON)

	// get response
	resp, err := e.execWithTracingAndLog(ctx, "Search Records", req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("an error happend while listing records (cause: %s)", err.Error()))
	}

	// http error handling
	if resp.StatusCode != http.StatusOK {
		return nil, e.handleInvalidHttpCode(*resp)
	}

	// build json result
	result := make([]*api.Record, 0)
	bodyBytes, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, errors.New(fmt.Sprintf("unable to parse the listed records: %s, %v", string(bodyBytes), err))
	}

	return result, nil
}

func (e *Client) CreateRecord(ctx context.Context, record *api.Record) (*api.Record, error) {
	b, err := json.Marshal(record)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("an error happend while building request body for creating record %s! (cause: %s)", record.Message, err.Error()))
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/records", e.config.ApiUrl), bytes.NewReader(b))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("an error happend while building request for creating record %s! (cause: %s)", record.Message, err.Error()))
	}

	req.Header.Set(headerNameAccept, headerValueApplicationJSON)

	// get response
	resp, err := e.execWithTracingAndLog(ctx, "CreateRecord", req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("an error happend while creating domain %s (cause: %s)", record.Message, err.Error()))
	}

	// http error handling
	if resp.StatusCode != http.StatusCreated {
		return nil, e.handleInvalidHttpCode(*resp)
	}

	// build json result
	result := &api.Record{}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if json.Unmarshal(bodyBytes, result) != nil {
		return nil, errors.New(fmt.Sprintf("unable to parse the created record: %s", string(bodyBytes)))
	}
	return result, nil
}
