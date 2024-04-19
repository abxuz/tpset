package mercury

import (
	"bytes"
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
	addr    string
	stok    string
	cookies []*http.Cookie
}

func NewService(addr string) *Service {
	return &Service{addr: addr}
}

func (s *Service) IsAC(ctx context.Context) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.addr+"/webpages/index.html", nil)
	if err != nil {
		return false, err
	}

	resp, err := http.DefaultClient.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return false, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	return bytes.Contains(data, []byte("mercury")), nil
}

func (s *Service) request(ctx context.Context, uri string, data string) ([]byte, error) {
	uri = s.addr + "/cgi-bin/luci/;stok=" + s.stok + uri
	params := make(url.Values)
	params.Set("data", data)
	body := strings.NewReader(params.Encode())
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Set("Referer", s.addr+"/webpages/index.html")
	for _, cookie := range s.cookies {
		request.AddCookie(cookie)
	}

	resp, err := http.DefaultClient.Do(request)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	retErr := errors.New(string(respData))

	var m Object
	err = json.Unmarshal(respData, &m)
	if err != nil {
		return nil, retErr
	}

	errno, ok := m["error_code"].(string)
	if !ok || errno != "0" {
		return nil, retErr
	}

	return respData, nil
}
