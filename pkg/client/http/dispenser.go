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

func (e *Client) GetDispenser(ctx context.Context, ref string) (*api.Dispenser, error) {
	// create request
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/dispensers/%s", e.config.ApiUrl, ref), nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("an error happend while building request for getting team %s! (cause: %s)", ref, err.Error()))
	}

	req.Header.Set(headerNameAccept, headerValueApplicationJSON)

	// get response
	resp, err := e.execWithTracingAndLog(ctx, "Get Dispenser", req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("an error happend while getting dispenser %s (cause: %s)", ref, err.Error()))
	}

	// http error handling
	if resp.StatusCode != http.StatusOK {
		return nil, e.handleInvalidHttpCode(*resp)
	}

	// build json result
	result := &api.Dispenser{}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if json.Unmarshal(bodyBytes, result) != nil {
		return nil, errors.New(fmt.Sprintf("unable to parse the gotten dispenser: %s", string(bodyBytes)))
	}

	return result, nil
}

func (e *Client) SearchDispensers(ctx context.Context) ([]*api.Dispenser, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/dispensers", e.config.ApiUrl), nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("an error happend while building request for listing dispensers (cause: %s)", err.Error()))
	}

	req.Header.Set(headerNameAccept, headerValueApplicationJSON)

	// get response
	resp, err := e.execWithTracingAndLog(ctx, "Search Dispensers", req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("an error happend while listing dispensers (cause: %s)", err.Error()))
	}

	// http error handling
	if resp.StatusCode != http.StatusOK {
		return nil, e.handleInvalidHttpCode(*resp)
	}

	// build json result
	result := make([]*api.Dispenser, 0)
	bodyBytes, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, errors.New(fmt.Sprintf("unable to parse the listed dispensers: %s, %v", string(bodyBytes), err))
	}

	return result, nil
}

func (e *Client) CreateDispenser(ctx context.Context, dispenser *api.DispenserEditable) (*api.Dispenser, error) {
	b, err := json.Marshal(dispenser)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("an error happend while building request body for creating dispenser %s! (cause: %s)", dispenser.Beer, err.Error()))
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/dispensers", e.config.ApiUrl), bytes.NewReader(b))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("an error happend while building request for creating dispenser %s! (cause: %s)", dispenser.Beer, err.Error()))
	}

	req.Header.Set(headerNameAccept, headerValueApplicationJSON)

	// get response
	resp, err := e.execWithTracingAndLog(ctx, "CreateDispenser", req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("an error happend while creating domain %s (cause: %s)", dispenser.Beer, err.Error()))
	}

	// http error handling
	if resp.StatusCode != http.StatusCreated {
		return nil, e.handleInvalidHttpCode(*resp)
	}

	// build json result
	result := &api.Dispenser{}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if json.Unmarshal(bodyBytes, result) != nil {
		return nil, errors.New(fmt.Sprintf("unable to parse the created dispenser: %s", string(bodyBytes)))
	}
	return result, nil
}

func (e *Client) UpdateDispenser(ctx context.Context, dispenser *api.Dispenser) (*api.Dispenser, error) {
	b, err := json.Marshal(dispenser)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("an error happend while building request body for saving dispenser %s! (cause: %s)", dispenser.Ref, err.Error()))
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/dispensers/%s", e.config.ApiUrl, dispenser.Ref), bytes.NewReader(b))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("an error happend while building request for updating dispenser %s! (cause: %s)", dispenser.Ref, err.Error()))
	}

	req.Header.Set(headerNameAccept, headerValueApplicationJSON)

	// get response
	resp, err := e.execWithTracingAndLog(ctx, "UpdateDispense", req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("an error happend while updating dispenser %s (cause: %s)", dispenser.Ref, err.Error()))
	}

	// http error handling
	if resp.StatusCode != http.StatusOK {
		return nil, e.handleInvalidHttpCode(*resp)
	}

	// build json result
	result := &api.Dispenser{}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if json.Unmarshal(bodyBytes, result) != nil {
		return nil, errors.New(fmt.Sprintf("unable to parse the updated dispenser: %s", string(bodyBytes)))
	}

	return result, nil
}

func (e *Client) DeleteDispenser(ctx context.Context, name string) error {

	// create request
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/dispensers/%s", e.config.ApiUrl, name), nil)
	if err != nil {
		return errors.New(fmt.Sprintf("an error happend while building request for deleting dispenser %s! (cause: %s)", name, err.Error()))
	}

	req.Header.Set(headerNameAccept, headerValueApplicationJSON)

	// get response
	resp, err := e.execWithTracingAndLog(ctx, "Delete dispenser", req)
	if err != nil {
		return errors.New(fmt.Sprintf("an error happend while deleting dispenser %s (cause: %s)", name, err.Error()))
	}

	// http error handling
	if resp.StatusCode != http.StatusNoContent {
		return e.handleInvalidHttpCode(*resp)
	}

	return nil
}
