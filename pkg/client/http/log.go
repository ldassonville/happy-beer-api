package http

import (
	"context"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strings"
	"time"
)

func (e *Client) execWithTracingAndLog(ctx context.Context, reqTitle string, req *http.Request) (*http.Response, error) {

	if e.config.DebugHttp {
		e.config.DebugLogger("HTTP %s", reqTitle)
		e.traceRequest(req)
		start := time.Now()
		resp, err := e.httpClient.Do(req)
		elapsed := time.Since(start)
		e.config.DebugLogger("[ HTTP RESPONSE TIME: %s ]", elapsed)
		if resp != nil {
			e.traceResponse(resp)
		}
		return resp, err
	} else {
		return e.httpClient.Do(req)
	}
}

func (e *Client) traceRequest(req *http.Request) {
	if e.config.DebugHttp {
		dump, err := httputil.DumpRequestOut(req, true)
		if err == nil {
			e.config.DebugLogger(dump2string(dump, "> "))
		} else {
			e.config.DebugLogger("ERROR: cannot trace request string: %s", err.Error())
		}
	}
}

// trace response only if debug mode
func (e *Client) traceResponse(resp *http.Response) {
	if e.config.DebugHttp {
		dump, err := httputil.DumpResponse(resp, true)
		if err == nil {
			e.config.DebugLogger(dump2string(dump, "< "))
		} else {
			e.config.DebugLogger("ERROR: cannot trace request string: %s", err.Error())
		}
	}
}

var (
	textBeginningRegex = regexp.MustCompile("^")
	newLineRegex       = regexp.MustCompile("\n")
)

func dump2string(dump []byte, prefix string) string {
	dumpStr := string(dump)
	dumpStr = strings.Trim(dumpStr, " \r\n\t")
	dumpStr = textBeginningRegex.ReplaceAllString(dumpStr, prefix)
	dumpStr = newLineRegex.ReplaceAllString(dumpStr, "\n"+prefix)
	return dumpStr
}
