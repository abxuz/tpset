package mercury

import (
	"context"
	"encoding/json"
	"errors"
)

type Group struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Seq     string `json:"seq"`
	Active  int    `json:"active"`
	Total   int    `json:"total"`
	TempNum int    `json:"temp_num"`
}

func (s *Service) ListGroup(ctx context.Context) ([]Group, error) {
	resp, err := s.request(ctx, "/admin/apmngr?form=grouplist", `{"method":"get","params":{}}`)
	if err != nil {
		return nil, err
	}

	retErr := errors.New(string(resp))

	var result struct {
		Result []Group `json:"result"`
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, retErr
	}

	return result.Result, nil
}
