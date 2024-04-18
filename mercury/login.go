package mercury

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func encrypt(msg string) string {
	e, _ := strconv.ParseInt("010001", 16, 0)
	n, _ := new(big.Int).SetString(`D1E79FF135D14E342D76185C23024E6DEAD4D6EC2C317A526C811E83538EA4E5ED8E1B0EEE5CE26E3C1B6A5F1FE11FA804F28B7E8821CA90AFA5B2F300DF99FDA27C9D2131E031EA11463C47944C05005EF4C1CE932D7F4A87C7563581D9F27F0C305023FCE94997EC7D790696E784357ED803A610EBB71B12A8BE5936429BFD`, 16)
	m := nopadding([]rune(msg), (n.BitLen()+7)>>3)
	c := new(big.Int).Exp(m, big.NewInt(e), n).Bytes()
	return hex.EncodeToString(c)
}

func nopadding(s []rune, n int) *big.Int {
	length := len(s)
	ba := make([]byte, n)
	i, j := 0, 0
	for i < length && j < n {
		c := s[i]
		i++
		if c < 128 {
			ba[j] = byte(c)
			j++
		} else if c > 127 && c < 2048 {
			ba[j] = byte((c & 63) | 128)
			j++
			ba[j] = byte((c >> 6) | 192)
			j++
		} else {
			ba[j] = byte((c & 63) | 128)
			j++
			ba[j] = byte(((c >> 6) & 63) | 128)
			j++
			ba[j] = byte((c >> 12) | 124)
		}
	}

	for j < n {
		ba[j] = 0
		j++
	}

	return new(big.Int).SetBytes(ba)
}

func (s *Service) Login(ctx context.Context, username, password string) error {
	params := make(url.Values)
	params.Set("data", `{"method":"login","params":{"username":"`+username+`","password":"`+encrypt(password)+`"}}`)

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
