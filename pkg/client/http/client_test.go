package http

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/ldassonville/beer-puller-api/pkg/client"
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
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return errors.New("unauthorized by API server")
	case http.StatusForbidden:
		return errors.New("forbidden by API server (you are probably missing access rights)")
	case http.StatusBadRequest:
		bodyBytes, _ := io.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf("invalid request to API server (probably due to some input value)\nbody: %s", string(bodyBytes)))
	case http.StatusNotFound:
		bodyBytes, _ := io.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf("not found in API server\nbody: %s", string(bodyBytes)))
	case http.StatusServiceUnavailable:
		return errors.New("API server currently down")
	case http.StatusPreconditionFailed:
		return errors.New("precondition failed")
	case http.StatusConflict:
		return errors.New("resource conflict")
	case http.StatusInternalServerError:
		bodyBytes, _ := io.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf("an error happened on the API server\nbody: %s", string(bodyBytes)))
	default:
		bodyBytes, _ := io.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf("an unknown error happened while requesting the API server\ncode: %d\nbody: %s", resp.StatusCode, string(bodyBytes)))
	}
}
