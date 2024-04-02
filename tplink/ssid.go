package tplink

import (
	"context"
	"encoding/json"
	"errors"
)

type SSID struct {
	SSIDId             string `json:"ssid_id"`
	ServId             string `json:"serv_id"`
	SSID               string `json:"ssid"`
	SSIDCodeType       string `json:"ssid_code_type"`
	Key                string `json:"key"`
	KeyUpdateIntv      string `json:"key_update_intv"`
	Enable             string `json:"enable"`
	Encryption         string `json:"encryption"`
	Cipher             string `json:"cipher"`
	Auth               string `json:"auth"`
	PrivKey            string `json:"priv_key"`
	Isolate            string `json:"isolate"`
	SSIDBrd            string `json:"ssidbrd"`
	UpLimit            string `json:"up_limit"`
	DownLimit          string `json:"down_limit"`
	IGMPQueryInterval  string `json:"igmp_query_interval"`
	IGMPMaxRespTime    string `json:"igmp_max_resp_time"`
	IGMPSnoopingSwitch string `json:"igmp_snooping_switch"`
	IGMPSnoopingEnable string `json:"igmp_snooping_enable"`
	IGMPDropUnknown    string `json:"igmp_drop_unknown"`
	IGMPQuerySwitch    string `json:"igmp_query_switch"`
	IGMP2gRate         string `json:"igmp_2g_rate"`
	IGMP5gRate         string `json:"igmp_5g_rate"`
	RadiusIP           string `json:"radius_ip"`
	RadiusPwd          string `json:"radius_pwd"`
	RadiusPort         string `json:"radius_port"`
	RadiusAcctPort     string `json:"radius_acct_port"`
	AutoBind           string `json:"auto_bind"`
	DefaultBindVlan    string `json:"default_bind_vlan"`
	DefaultBindFreq    string `json:"default_bind_freq"`
	BwCtrlEnable       string `json:"bw_ctrl_enable"`
	BwCtrlMode         string `json:"bw_ctrl_mode"`
	Desc               string `json:"desc"`
}

func (s *SSID) escapeList() []*string {
	return []*string{&s.SSID, &s.Desc}
}

func (s *Service) ListSSID(ctx context.Context) ([]SSID, error) {
	body := `{"method":"get","apmng_wserv":{"table":"wlan_serv","para":{"start":0,"end":1999}}}`
	data, err := s.ds(ctx, body)
	if err != nil {
		return nil, err
	}
	retErr := errors.New(string(data))

	var result struct {
		ApMngWServ struct {
			WlanServ []map[string]json.RawMessage `json:"wlan_serv"`
		} `json:"apmng_wserv"`
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, retErr
	}

	list := make([]SSID, 0)
	for _, v := range result.ApMngWServ.WlanServ {
		if len(v) != 1 {
			return nil, retErr
		}

		var raw json.RawMessage
		for _, raw = range v {
			break
		}

		var ssid SSID
		err = json.Unmarshal(raw, &ssid)
		if err != nil {
			return nil, retErr
		}
		err = unescape(&ssid)
		if err != nil {
			return nil, err
		}
		list = append(list, ssid)
	}

	return list, nil
}

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
