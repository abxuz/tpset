package mercury

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
)

type SSID struct {
	Id             string `json:"id"`
	CloudSSIDId    string `json:"cloud_ssid_id"`
	SSID           string `json:"ssid"`
	Encode         string `json:"encode"`
	Psk            string `json:"psk"`
	KeyUptIntv     string `json:"key_upt_intv"`
	Enable         string `json:"enable"`
	AuthType       string `json:"auth_type"`
	SecType        string `json:"sec_type"`
	EncryptAlg     string `json:"encrypt_alg"`
	PrivKey        string `json:"priv_key"`
	Isolation      string `json:"isolation"`
	Hide           string `json:"hide"`
	UpstreamRate   string `json:"upstream_rate"`
	DownstreamRate string `json:"downstream_rate"`
	RadiusServ     string `json:"radius_serv"`
	RadiusPwd      string `json:"radius_pwd"`
	RadiusPort     string `json:"radius_port"`
	RadiusAcctPort string `json:"radius_acct_port"`
	BwEnable       string `json:"bw_enable"`
	BwMode         string `json:"bw_mode"`
	Descr          string `json:"descr"`
}

func (s *Service) ListSSID(ctx context.Context) ([]SSID, error) {
	uri := "/admin/ac_wservice?form=wserv_list"
	_, err := s.request(ctx, uri, `{"method":"change","params":{"pageSize":"100"}}`)
	if err != nil {
		return nil, err
	}

	list := make([]SSID, 0)
	for i := 0; true; i++ {
		data, err := s.request(ctx, uri, `{"method":"get","params":{"pageNo":`+strconv.Itoa(i)+`}}`)
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

		var arr []SSID
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

/*
func (s *Service) AddSSID(ctx context.Context, ssid SSID) error {
	escape(&ssid)
	if ssid.UpLimit == "0" || ssid.UpLimit == "" {
		ssid.UpLimit = "128"
	}
	if ssid.DownLimit == "0" || ssid.DownLimit == "" {
		ssid.DownLimit = "128"
	}
	body := `{
	"method":"add",
	"apmng_wserv":{
		"table":"wlan_serv",
		"para":{
			"serv_id":"` + ssid.ServId + `",
			"enable":"` + ssid.Enable + `",
			"ssid":"` + ssid.SSID + `",
			"ssid_code_type":"` + ssid.SSIDCodeType + `",
			"desc":"` + ssid.Desc + `",
			"isolate":"` + ssid.Isolate + `",
			"ssidbrd":"` + ssid.SSIDBrd + `",
			"encryption":"` + ssid.Encryption + `",
			"auth":"` + ssid.Auth + `",
			"cipher":"` + ssid.Cipher + `",
			"key_update_intv":"` + ssid.KeyUpdateIntv + `",
			"key":"` + ssid.Key + `",
			"bw_ctrl_enable":"` + ssid.BwCtrlEnable + `",
			"bw_ctrl_mode":"` + ssid.BwCtrlMode + `",
			"up_limit":"` + ssid.UpLimit + `",
			"down_limit":"` + ssid.DownLimit + `",
			"auto_bind":"` + ssid.AutoBind + `",
			"default_bind_freq":"` + ssid.DefaultBindFreq + `",
			"default_bind_vlan":"` + ssid.DefaultBindVlan + `"
		}
	}}`
	_, err := s.ds(ctx, body)
	return err
}

func (s *Service) SetSSID(ctx context.Context, ssid SSID) error {
	escape(&ssid)
	if ssid.UpLimit == "0" || ssid.UpLimit == "" {
		ssid.UpLimit = "128"
	}
	if ssid.DownLimit == "0" || ssid.DownLimit == "" {
		ssid.DownLimit = "128"
	}
	body := `{
	"method":"set",
	"apmng_wserv":{
		"table":"wlan_serv",
		"filter":[{"serv_id":"` + ssid.ServId + `"}],
		"para":{
			"serv_id":"` + ssid.ServId + `",
			"enable":"` + ssid.Enable + `",
			"ssid":"` + ssid.SSID + `",
			"ssid_code_type":"` + ssid.SSIDCodeType + `",
			"desc":"` + ssid.Desc + `",
			"isolate":"` + ssid.Isolate + `",
			"ssidbrd":"` + ssid.SSIDBrd + `",
			"encryption":"` + ssid.Encryption + `",
			"auth":"` + ssid.Auth + `",
			"cipher":"` + ssid.Cipher + `",
			"key_update_intv":"` + ssid.KeyUpdateIntv + `",
			"key":"` + ssid.Key + `",
			"bw_ctrl_enable":"` + ssid.BwCtrlEnable + `",
			"bw_ctrl_mode":"` + ssid.BwCtrlMode + `",
			"up_limit":"` + ssid.UpLimit + `",
			"down_limit":"` + ssid.DownLimit + `",
			"auto_bind":"` + ssid.AutoBind + `",
			"default_bind_freq":"` + ssid.DefaultBindFreq + `",
			"default_bind_vlan":"` + ssid.DefaultBindVlan + `"
		}
	}}`
	_, err := s.ds(ctx, body)
	return err
}
*/
