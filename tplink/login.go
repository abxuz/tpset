package tplink

import (
	"context"
	"encoding/json"
	"errors"
)

func (s *Service) Login(ctx context.Context, username, password string) error {
	body := `{
	"method":"do",
	"login":{
		"username":"` + username + `",
		"password":"` + password + `"
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
