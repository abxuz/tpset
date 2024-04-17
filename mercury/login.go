package mercury

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (s *Service) Login(ctx context.Context, username, password string) error {
	params := make(url.Values)
	params.Set("data", `{"method":"login","params":{"username":"`+username+`","password":"`+password+`"}}`)

	uri := s.addr + "/cgi-bin/luci/;stok=/login?form=login"
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Set("Referer", s.addr+"/webpages/login.html")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	respData, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}
	retErr := errors.New(string(respData))

	var m struct {
		Result *struct {
			Stok string `json:"stok"`
		} `json:"result"`
	}
	err = json.Unmarshal(respData, &m)
	if err != nil {
		return retErr
	}

	if m.Result == nil || m.Result.Stok == "" {
		return retErr
	}

	s.stok = m.Result.Stok
	s.cookies = resp.Cookies()
	return nil
}
