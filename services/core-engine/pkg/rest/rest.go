package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gerins/log"
)

const (
	ProcessIDContextKey = "processID"
	ContentType         = "Content-Type"
	ApplicationJSON     = "application/json"
)

type rest struct {
	client *http.Client
}

type Rest interface {
	Post(ctx context.Context, url string, header map[string]string, payload interface{}) ([]byte, int, error)
	Put(ctx context.Context, url string, header map[string]string, payload interface{}) ([]byte, int, error)
	Get(ctx context.Context, url string, header map[string]string, queryParams map[string]string) ([]byte, int, error)
	Delete(ctx context.Context, url string, header map[string]string, queryParams map[string]string) ([]byte, int, error)
}

func New(timeout time.Duration) Rest {
	return &rest{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (r *rest) Post(ctx context.Context, url string, header map[string]string, payload interface{}) ([]byte, int, error) {
	var (
		req             *http.Request
		resp            *http.Response
		response        []byte
		err             error
		requestDuration int64        // Total duration when making request
		startTime       = time.Now() // Time when making request
	)

	defer func() {
		if resp != nil {
			log.Tracing(log.Context(ctx).ProcessID(), url, "POST", resp.StatusCode, response, header, payload, resp.Header, err, requestDuration)
		} else {
			log.Tracing(log.Context(ctx).ProcessID(), url, "POST", 0, response, header, payload, nil, err, requestDuration)
		}
	}()

	// Convert payload to []byte type
	requestPayload, err := json.Marshal(payload)
	if err != nil {
		log.Context(ctx).Errorf("error marshaling payload, %v", err)
		return nil, 0, err
	}

	// Creating new request with context
	req, err = http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestPayload))
	if err != nil {
		log.Context(ctx).Errorf("error creating new request, %v", err)
		return nil, 0, err
	}

	// Add ProcessID Header, super useful for tracing Log if we ecounter issue in another Service
	header[ProcessIDContextKey] = log.Context(ctx).ProcessID()
	header[ContentType] = ApplicationJSON

	// Adding header to the request
	for key, value := range header {
		req.Header.Set(key, value)
	}

	// Execute http request
	resp, err = r.client.Do(req)
	if err != nil {
		log.Context(ctx).Errorf("error when making request, %v", err)
		return nil, 0, err
	}
	defer resp.Body.Close()

	// Get the response body
	response, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Context(ctx).Errorf("error reading request body, %v", err)
		return nil, resp.StatusCode, err
	}

	requestDuration = time.Since(startTime).Milliseconds()
	return response, resp.StatusCode, nil
}

func (r *rest) Put(ctx context.Context, url string, header map[string]string, payload interface{}) ([]byte, int, error) {
	var (
		req             *http.Request
		resp            *http.Response
		response        []byte
		err             error
		requestDuration int64        // Total duration when making request
		startTime       = time.Now() // Time when making request
	)

	defer func() {
		if resp != nil {
			log.Tracing(log.Context(ctx).ProcessID(), url, "PUT", resp.StatusCode, response, header, payload, resp.Header, err, requestDuration)
		} else {
			log.Tracing(log.Context(ctx).ProcessID(), url, "PUT", 0, response, header, payload, nil, err, requestDuration)
		}
	}()

	// Convert payload to []byte type
	requestPayload, err := json.Marshal(payload)
	if err != nil {
		log.Context(ctx).Errorf("error marshaling payload, %v", err)
		return nil, 0, err
	}

	// Creating new request with context
	req, err = http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(requestPayload))
	if err != nil {
		log.Context(ctx).Errorf("error creating new request, %v", err)
		return nil, 0, err
	}

	// Add ProcessID Header, super useful for tracing Log if we ecounter issue in another Service
	header[ProcessIDContextKey] = log.Context(ctx).ProcessID()
	header[ContentType] = ApplicationJSON

	// Adding header to the request
	for key, value := range header {
		req.Header.Set(key, value)
	}

	// Execute http request
	resp, err = r.client.Do(req)
	if err != nil {
		log.Context(ctx).Errorf("error when making request, %v", err)
		return nil, 0, err
	}
	defer resp.Body.Close()

	// Get the response body
	response, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Context(ctx).Errorf("error reading request body, %v", err)
		return nil, resp.StatusCode, err
	}

	requestDuration = time.Since(startTime).Milliseconds()
	return response, resp.StatusCode, nil
}

func (r *rest) Get(ctx context.Context, url string, header map[string]string, queryParams map[string]string) ([]byte, int, error) {
	var (
		req             *http.Request
		resp            *http.Response
		response        []byte
		err             error
		requestDuration int64        // Total duration when making request
		startTime       = time.Now() // Time when making request
	)

	defer func() {
		if resp != nil {
			log.Tracing(log.Context(ctx).ProcessID(), url, "GET", resp.StatusCode, response, header, queryParams, resp.Header, err, requestDuration)
		} else {
			log.Tracing(log.Context(ctx).ProcessID(), url, "GET", 0, response, header, queryParams, nil, err, requestDuration)
		}
	}()

	// Creating new request with context
	req, err = http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Context(ctx).Errorf("error creating new request, %v", err)
		return nil, 0, err
	}

	// Add ProcessID Header, super useful for tracing Log if we ecounter issue in another Service
	header[ProcessIDContextKey] = log.Context(ctx).ProcessID()
	header[ContentType] = ApplicationJSON

	// Adding header to the request
	for key, value := range header {
		req.Header.Set(key, value)
	}

	// Building query params
	query := req.URL.Query()
	for key, value := range queryParams {
		query.Add(key, value)
	}

	// Add query params to the url
	req.URL.RawQuery = query.Encode()

	// Execute http request
	resp, err = r.client.Do(req)
	if err != nil {
		log.Context(ctx).Errorf("error when making request, %v", err)
		return nil, 0, err
	}
	defer resp.Body.Close()

	// Get the response body
	response, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Context(ctx).Errorf("error reading request body, %v", err)
		return nil, resp.StatusCode, err
	}

	requestDuration = time.Since(startTime).Milliseconds()
	return response, resp.StatusCode, nil
}

func (r *rest) Delete(ctx context.Context, url string, header map[string]string, queryParams map[string]string) ([]byte, int, error) {
	var (
		req             *http.Request
		resp            *http.Response
		response        []byte
		err             error
		requestDuration int64        // Total duration when making request
		startTime       = time.Now() // Time when making request
	)

	defer func() {
		if resp != nil {
			log.Tracing(log.Context(ctx).ProcessID(), url, "DELETE", resp.StatusCode, response, header, queryParams, resp.Header, err, requestDuration)
		} else {
			log.Tracing(log.Context(ctx).ProcessID(), url, "DELETE", 0, response, header, queryParams, nil, err, requestDuration)
		}
	}()

	// Creating new request with context
	req, err = http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		log.Context(ctx).Errorf("error creating new request, %v", err)
		return nil, 0, err
	}

	// Add ProcessID Header, super useful for tracing Log if we ecounter issue in another Service
	header[ProcessIDContextKey] = log.Context(ctx).ProcessID()
	header[ContentType] = ApplicationJSON

	// Adding header to the request
	for key, value := range header {
		req.Header.Set(key, value)
	}

	// Building query params
	query := req.URL.Query()
	for key, value := range queryParams {
		query.Add(key, value)
	}

	// Add query params to the url
	req.URL.RawQuery = query.Encode()

	// Execute http request
	resp, err = r.client.Do(req)
	if err != nil {
		log.Context(ctx).Errorf("error when making request, %v", err)
		return nil, 0, err
	}
	defer resp.Body.Close()

	// Get the response body
	response, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Context(ctx).Errorf("error reading request body, %v", err)
		return nil, resp.StatusCode, err
	}

	requestDuration = time.Since(startTime).Milliseconds()
	return response, resp.StatusCode, nil
}
