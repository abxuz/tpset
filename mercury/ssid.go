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

func (s *Service) AddSSID(ctx context.Context, ssid SSID) error {
	if ssid.Enable != "" && ssid.Enable != "0" {
		ssid.Enable = "1"
	} else {
		ssid.Enable = "0"
	}

	if ssid.Isolation != "" && ssid.Isolation != "0" {
		ssid.Isolation = "1"
	} else {
		ssid.Isolation = "0"
	}

	if ssid.Hide != "" && ssid.Hide != "0" {
		ssid.Hide = "1"
	} else {
		ssid.Hide = "0"
	}

	if ssid.SecType == "" {
		ssid.SecType = "0"
	}

	data := `{
    "method": "add",
    "params": {
        "index": 0,
        "old": "add",
        "new": {
            "id": "",
            "enable": "` + ssid.Enable + `",
            "ssid": "` + ssid.SSID + `",
            "encode": "1",
            "descr": "` + ssid.Descr + `",
            "isolation": "` + ssid.Isolation + `",
            "hide": "` + ssid.Hide + `",
            "sec_type": "` + ssid.SecType + `",
            "auth_type": "3",
            "encrypt_alg": "0",
            "key_upt_intv": "86400",
            "psk": "` + ssid.Psk + `",
            "bw_enable": "0",
            "bw_mode": "0",
            "upstream_rate": "128",
            "downstream_rate": "128"
        },
        "id": "add"
    }
}`
	_, err := s.request(ctx, "/admin/ac_wservice?form=wserv_list", data)
	return err
}

func (s *Service) SetSSID(ctx context.Context, prev, cur SSID) error {
	if cur.UpstreamRate == "" || cur.UpstreamRate == "0" {
		cur.UpstreamRate = "128"
	}

	if cur.DownstreamRate == "" || cur.DownstreamRate == "0" {
		cur.DownstreamRate = "128"
	}

	data := `{
	"method": "set",
	"params": {
		"index": 0,
		"old": {
			"enable": "` + prev.Enable + `",
			"ssid": "` + prev.SSID + `",
			"encode": "` + prev.Encode + `",
			"descr": "` + prev.Descr + `",
			"isolation": "` + prev.Isolation + `",
			"hide": "` + prev.Hide + `",
			"sec_type": "` + prev.SecType + `",
			"auth_type": "` + prev.AuthType + `",
			"encrypt_alg": "` + prev.EncryptAlg + `",
			"radius_serv": "` + prev.RadiusServ + `",
			"radius_port": "` + prev.RadiusPort + `",
			"radius_acct_port": "` + prev.RadiusAcctPort + `",
			"radius_pwd": "` + prev.RadiusPwd + `",
			"key_upt_intv": "` + prev.KeyUptIntv + `",
			"psk": "` + prev.Psk + `",
			"bw_enable": "` + prev.BwEnable + `",
			"bw_mode": "` + prev.BwMode + `",
			"upstream_rate": "` + prev.UpstreamRate + `",
			"downstream_rate": "` + prev.DownstreamRate + `"
		},
		"new": {
			"id": "` + cur.Id + `",
			"enable": "` + cur.Enable + `",
			"ssid": "` + cur.SSID + `",
			"encode": "` + cur.Encode + `",
			"descr": "` + cur.Descr + `",
			"isolation": "` + cur.Isolation + `",
			"hide": "` + cur.Hide + `",
			"sec_type": "` + cur.SecType + `",
			"auth_type": "` + cur.AuthType + `",
			"encrypt_alg": "` + cur.EncryptAlg + `",
			"radius_serv": "` + cur.RadiusServ + `",
			"radius_port": "` + cur.RadiusPort + `",
			"radius_acct_port": "` + cur.RadiusAcctPort + `",
			"radius_pwd": "` + cur.RadiusPwd + `",
			"key_upt_intv": "` + cur.KeyUptIntv + `",
			"psk": "` + cur.Psk + `",
			"bw_enable": "` + cur.BwEnable + `",
			"bw_mode": "` + cur.BwMode + `",
			"upstream_rate": "` + cur.UpstreamRate + `",
			"downstream_rate": "` + cur.DownstreamRate + `"
		},
		"id": "` + prev.Id + `"
	}
}`
	_, err := s.request(ctx, "/admin/ac_wservice?form=wserv_list", data)
	return err
}
