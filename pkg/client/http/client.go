package http

import (
	"crypto/tls"
	"fmt"
	"github.com/ldassonville/happy-beer-api/pkg/client"
	"io"
	"net/http"
	"time"
)

const (
	headerNameAccept = "Accept"

	headerValueApplicationJSON = "application/json"
)

type DebugLogger func(message string, args ...interface{})

type Config struct {
	ApiUrl             string
	DebugHttp          bool
	DebugLogger        DebugLogger
	InsecureSkipVerify bool
}

type Client struct {
	httpClient *http.Client
	config     *Config
}

func NewClient(config *Config) client.Client {

	var tr http.RoundTripper = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.InsecureSkipVerify},
	}

	return &Client{
		httpClient: &http.Client{
			Timeout:   time.Second * 10,
			Transport: tr,
		},
		config: config,
	}
}

func (e *Client) handleInvalidHttpCode(resp http.Response) error {

	res := &client.APIError{
		Code: resp.StatusCode,
	}
	bodyBytes, _ := io.ReadAll(resp.Body)

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		res.Reason = client.StatusReasonUnauthorized
		res.Message = "unauthorized by API server (you are probably missing authentication)"

	case http.StatusForbidden:
		res.Reason = client.StatusReasonForbidden
		res.Message = "forbidden by API server (you are probably missing access rights)"

	case http.StatusBadRequest:
		res.Reason = client.StatusBadRequest
		res.Message = fmt.Sprintf("invalid request to API server (probably due to some input value)\nbody: %s", string(bodyBytes))

	case http.StatusNotFound:
		res.Reason = client.StatusReasonNotFound
		res.Message = "resource not found"

	case http.StatusServiceUnavailable:
		res.Reason = client.StatusServiceUnavailable
		res.Message = "API server currently down"

	case http.StatusPreconditionFailed:
		res.Reason = client.StatusPreconditionFailed
		res.Message = "precondition failed"

	case http.StatusConflict:
		res.Reason = client.StatusConflict
		res.Message = "resource conflict"

	case http.StatusInternalServerError:
		res.Reason = client.StatusInternalError
		res.Message = fmt.Sprintf("an error happened on the API server\nbody: %s", string(bodyBytes))

	default:
		res.Reason = client.StatusReasonUnknown
		res.Message = fmt.Sprintf("an unknown error happened while requesting the API server\ncode: %d\nbody: %s", resp.StatusCode, string(bodyBytes))
	}
	return res
}
