package mercury

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
)

type AP struct {
	Id                string `json:"id"`
	EntryName         string `json:"entry_name"`
	EntryType         string `json:"entry_type"`
	GroupId           string `json:"group_id"`
	ModelId           string `json:"model_id"`
	CloudConfigId     string `json:"cloud_config_id"`
	Seq               string `json:"seq"`
	HardwareVer       string `json:"hardware_ver"`
	SoftwareVer       string `json:"software_ver"`
	Mac               string `json:"mac"`
	ApTrunk           string `json:"ap_trunk"`
	Status            string `json:"status"`
	Errno             string `json:"errno"`
	Led               string `json:"led"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	ApKeepAlive       string `json:"ap_keep_alive"`
	ClientKeepAlive   string `json:"client_keep_alive"`
	ClientIdle        string `json:"client_idle"`
	PhyWireVlan1      string `json:"phy_wire_vlan_1"`
	PhyWireVlan2      string `json:"phy_wire_vlan_2"`
	PhyWireVlan3      string `json:"phy_wire_vlan_3"`
	PhyWireVlan4      string `json:"phy_wire_vlan_4"`
	OfflineSelfmanage string `json:"offline_selfmanage"`
}

func (s *Service) ListAP(ctx context.Context, gid string) ([]AP, error) {
	uri := "/admin/apmngr?form=aplist"
	_, err := s.request(ctx, uri, `{"method":"change","params":{"pageSize":"500"}}`)
	if err != nil {
		return nil, err
	}

	list := make([]AP, 0)
	for i := 0; true; i++ {
		data, err := s.request(ctx, uri, `{"method":"get","params":{"pageNo":`+strconv.Itoa(i)+`,"group_id":`+gid+`}}`)
		if err != nil {
			return nil, err
		}
		retErr := errors.New(string(data))

		var result struct {
			Others struct {
				TotalNum int `json:"totalNum"`
			} `json:"others"`
			Result json.RawMessage `json:"result"`
		}
		err = json.Unmarshal(data, &result)
		if err != nil {
			return nil, retErr
		}

		if result.Others.TotalNum == 0 {
			break
		}

		var arr []AP
		err = json.Unmarshal(result.Result, &arr)
		if err != nil {
			return nil, retErr
		}

		list = append(list, arr...)

		if len(list) == result.Others.TotalNum {
			break
		}
	}

	return list, nil
}

func (s *Service) SetAP(ctx context.Context, prev, cur AP) error {
	data :=
		`{
    "method": "set",
    "params": {
        "index": 0,
        "old": {
            "id": "` + prev.Id + `",
            "group_id": "` + prev.GroupId + `",
            "entry_name": "` + prev.EntryName + `",
            "entry_type": "` + prev.EntryType + `",
            "model_id": "` + prev.ModelId + `",
            "hardware_ver": "` + prev.HardwareVer + `",
            "software_ver": "` + prev.SoftwareVer + `",
            "mac": "` + prev.Mac + `",
            "ap_keep_alive": "` + prev.ApKeepAlive + `",
            "client_keep_alive": "` + prev.ClientKeepAlive + `",
            "client_idle": "` + prev.ClientIdle + `",
            "phy_wire_vlan_1": "` + prev.PhyWireVlan1 + `",
            "phy_wire_vlan_2": "` + prev.PhyWireVlan2 + `",
            "phy_wire_vlan_3": "` + prev.PhyWireVlan3 + `",
            "phy_wire_vlan_4": "` + prev.PhyWireVlan4 + `",
            "offline_selfmanage": "` + prev.OfflineSelfmanage + `",
            "username": "` + prev.Username + `",
            "password": "` + prev.Password + `",
            "ap_trunk": "` + prev.ApTrunk + `",
            "led": "` + prev.Led + `",
            "status": "` + prev.Status + `",
            "errno": "` + prev.Errno + `"
        },
        "new": {
            "id": "` + cur.Id + `",
            "group_id": "` + cur.GroupId + `",
            "entry_name": "` + cur.EntryName + `",
            "model_id": "` + cur.ModelId + `",
            "hardware_ver": "` + cur.HardwareVer + `",
            "entry_type": "` + cur.EntryType + `",
            "ap_keep_alive": "` + cur.ApKeepAlive + `",
            "client_keep_alive": "` + cur.ClientKeepAlive + `",
            "client_idle": "` + cur.ClientIdle + `",
            "offline_selfmanage": "` + cur.OfflineSelfmanage + `",
            "username": "` + cur.Username + `",
            "password": "` + cur.Password + `",
            "led": "` + cur.Led + `",
            "status": "` + cur.Seq + `",
            "hdap_peer_id": ""
        },
        "key": "key-0"
    }
}`
	_, err := s.request(ctx, "/admin/apmngr?form=aplist", data)
	return err
}
