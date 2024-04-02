package tplink

import (
	"context"
	"encoding/json"
	"errors"
)

type Group struct {
	GroupId string `json:"group_id"`
	Name    string `json:"name"`
	Seq     string `json:"seq"`
	Active  string `json:"active"`
	Total   string `json:"total"`
	TempNum string `json:"temp_num"`
}

func (g *Group) escapeList() []*string {
	return []*string{&g.Name}
}

func (s *Service) ListGroup(ctx context.Context) ([]Group, error) {
	body := `{"method":"get","apmng_set":{"table":"group_list"}}`
	data, err := s.ds(ctx, body)
	if err != nil {
		return nil, err
	}
	retErr := errors.New(string(data))

	var result struct {
		ApMngSet struct {
			GroupList []map[string]json.RawMessage `json:"group_list"`
		} `json:"apmng_set"`
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, retErr
	}

	list := make([]Group, 0)
	for _, v := range result.ApMngSet.GroupList {
		if len(v) != 1 {
			return nil, retErr
		}

		var raw json.RawMessage
		for _, raw = range v {
			break
		}

		var g Group
		err = json.Unmarshal(raw, &g)
		if err != nil {
			return nil, retErr
		}
		err = unescape(&g)
		if err != nil {
			return nil, err
		}
		list = append(list, g)
	}

	return list, nil
}
