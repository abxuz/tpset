package tplink

import (
	"context"
	"encoding/json"
)

type Bind struct {
	Unbind  bool
	ServId  string
	VlanId  string
	RadioId []string
}

func (s *Service) Bind(ctx context.Context, bind Bind) error {
	vlans := make([]string, len(bind.RadioId))
	for i := range vlans {
		vlans[i] = bind.VlanId
	}
	vlanId, err := json.Marshal(vlans)
	if err != nil {
		return err
	}
	radioId, err := json.Marshal(bind.RadioId)
	if err != nil {
		return err
	}

	status := "1"
	if bind.Unbind {
		status = "0"
	}

	body := `{
	"method":"set",
	"apmng_wserv":{
		"table":"radio_bind_list",
		"filter":[{"serv_id":"` + bind.ServId + `"}],
		"para":{
			"bind_status":"` + status + `",
			"radio_id":` + string(radioId) + `,
			"vlan_id":` + string(vlanId) + `
		}
	}}`
	_, err = s.ds(ctx, body)
	return err
}
