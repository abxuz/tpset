package tplink

import (
	"context"
	"encoding/json"
	"errors"
)

func encrypt(password string) string {
	var (
		a = []rune(password)
		b = []rune("RDpbLfCPsJZ7fiv")
		d = []rune("yLwVl0zKqws7LgKPRQ84Mdt708T1qQ3Ha7xv3H7NyU84p21BriUWBU43odz3iP4rBL3cD02KZciXTysVXiV8ngg6vL48rPJyAUw0HurW20xqxv9aYb4M9wK1Ae0wlro510qXeU07kV57fQMc8L6aLgMLwygtc0F10a0Dg70TOoouyFhdysuRMO51yY5ZlOZZLEal1h0t9YQW0Ko7oBwmCAHoic4HYbUyVeU3sfQ1xtXcPcf1aT303wAQhv66qzW")
	)

	e := []rune{}
	g := len(a)
	m := len(b)
	k := len(d)
	var (
		t rune
		l rune
	)

	var h int
	if g > m {
		h = g
	} else {
		h = m
	}

	for f := 0; f < h; f++ {
		t = 187
		l = 187
		if f >= g {
			t = b[f]
		} else {
			if f >= m {
				l = a[f]
			} else {
				l = a[f]
				t = b[f]
			}
		}
		e = append(e, d[int(l)^int(t)%int(k)])
	}
	return string(e)
}

func (s *Service) Login(ctx context.Context, username, password string) error {
	body := `{
	"method":"do",
	"login":{
		"username":"` + username + `",
		"password":"` + encrypt(password) + `"
	}}`
	data, err := s.request(ctx, s.addr+"/", body)
	if err != nil {
		return err
	}
	retErr := errors.New(string(data))

	var m Object
	err = json.Unmarshal(data, &m)
	if err != nil {
		return retErr
	}

	stok, ok := m["stok"].(string)
	if !ok || stok == "" {
		return retErr
	}

	s.stok = stok
	return nil
}
