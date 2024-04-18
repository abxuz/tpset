package mercury

import (
	"context"
	"strings"
)

type Bind struct {
	Unbind  bool
	ServId  string
	VlanId  string
	RadioId []string
}

func (s *Service) Bind(ctx context.Context, bind Bind) error {
	action := "bind"
	if bind.Unbind {
		action = "unbind"
	}

	ids := strings.Join(bind.RadioId, ",")
	data := `{
	"method": "set",
	"params": {
		"action": "` + action + `",
		"vlan_id": "` + bind.VlanId + `",
		"wserv_id": "` + bind.ServId + `",
		"rd_id": "` + ids + `"
	},
	"id": "` + ids + `"
}
`
	_, err := s.request(ctx, "/admin/ac_wservice?form=rd_list", data)
	return err
}
