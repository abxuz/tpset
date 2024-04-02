package tplink

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Object = map[string]any
type Array = []any

type Service struct {
	addr string
	stok string
}

func NewService(addr string) *Service {
	return &Service{addr: addr}
}

func (s *Service) ds(ctx context.Context, body string) ([]byte, error) {
	return s.request(ctx, s.addr+"/stok="+s.stok+"/ds", body)
}

func (s *Service) request(ctx context.Context, uri string, body string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		ErrorCode int `json:"error_code"`
	}
	result.ErrorCode = -1
	err = json.Unmarshal(data, &result)
	if err != nil || result.ErrorCode != 0 {
		return nil, errors.New(string(data))
	}
	return data, nil
}

type escaper interface {
	escapeList() []*string
}

func unescape(escaper escaper) (err error) {
	for _, p := range escaper.escapeList() {
		*p, err = url.QueryUnescape(*p)
		if err != nil {
			return
		}
	}
	return
}

func escape(escaper escaper) {
	for _, p := range escaper.escapeList() {
		*p = url.QueryEscape(*p)
	}
}
