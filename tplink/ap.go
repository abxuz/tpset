package tplink

import (
	"context"
	"encoding/json"
	"errors"
)

type AP struct {
	EntryId            string `json:"entry_id"`
	EntryName          string `json:"entry_name"`
	EntryType          string `json:"entry_type"`
	GroupId            string `json:"group_id"`
	ModelId            string `json:"model_id"`
	CapId              string `json:"cap_id"`
	CloudConfigId      string `json:"cloud_config_id"`
	Seq                string `json:"seq"`
	HardwareVer        string `json:"hardware_ver"`
	Mac                string `json:"mac"`
	ApTrunk            string `json:"ap_trunk"`
	NewApTrunk         string `json:"new_ap_trunk"`
	ApRole             string `json:"ap_role"`
	Status             string `json:"status"`
	TelnetStatus       string `json:"telnet_status"`
	Errno              string `json:"errno"`
	LagMode            string `json:"lag_mode"`
	Led                string `json:"led"`
	BluetoothEnable    string `json:"bluetooth_enable"`
	LabPhyNum          string `json:"lag_phy_num"`
	LagPhyList         string `json:"lag_phy_list"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	AccountUsername    string `json:"account_username"`
	AccountPassword    string `json:"account_password"`
	ApKeepAlive        string `json:"ap_keep_alive"`
	ClientKeepAlive    string `json:"client_keep_alive"`
	ClientIdle         string `json:"client_idle"`
	LedOpenDate        string `json:"led_open_date"`
	LedOpenTime        string `json:"led_open_time"`
	LedTimerEnable     string `json:"led_timer_enable"`
	LedWifiSync        string `json:"led_wifi_sync"`
	LedCloseDate       string `json:"led_close_date"`
	LedCloseTime       string `json:"led_close_time"`
	PhyWireVlan1       string `json:"phy_wire_vlan_1"`
	PhyWireVlan2       string `json:"phy_wire_vlan_2"`
	PhyWireVlan3       string `json:"phy_wire_vlan_3"`
	PhyWireVlan4       string `json:"phy_wire_vlan_4"`
	FitFatButtonStatus string `json:"fit_fat_button_status"`
	OfflineSelfmanage  string `json:"offline_selfmanage"`
}

func (a *AP) escapeList() []*string {
	return []*string{&a.LedCloseTime, &a.LedOpenTime}
}

func (s *Service) ListAP(ctx context.Context, gid string) ([]AP, error) {
	body := `{
	"method":"get",
	"apmng_set": {
		"table":"ap_list",
		"filter":[{"group_id":"` + gid + `"}],
		"para":{"start":0,"end":1999}
	}}`
	data, err := s.ds(ctx, body)
	if err != nil {
		return nil, err
	}
	retErr := errors.New(string(data))

	var result struct {
		ApMngSet struct {
			ApList []map[string]json.RawMessage `json:"ap_list"`
		} `json:"apmng_set"`
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, retErr
	}

	list := make([]AP, 0)
	for _, v := range result.ApMngSet.ApList {
		if len(v) != 1 {
			return nil, retErr
		}

		var raw json.RawMessage
		for _, raw = range v {
			break
		}

		var ap AP
		err = json.Unmarshal(raw, &ap)
		if err != nil {
			return nil, retErr
		}
		err = unescape(&ap)
		if err != nil {
			return nil, err
		}
		list = append(list, ap)
	}

	return list, nil
}

func (s *Service) SetAP(ctx context.Context, ap AP) error {
	escape(&ap)
	body := `{
	"method":"set",
	"apmng_set":{
		"table":"ap_list",
		"filter":[{"entry_id":"` + ap.EntryId + `"}],
		"para":{
			"group_id":"` + ap.GroupId + `",
			"entry_name":"` + ap.EntryName + `",
			"ap_keep_alive":"` + ap.ApKeepAlive + `",
			"client_keep_alive":"` + ap.ClientKeepAlive + `",
			"client_idle":"` + ap.ClientIdle + `",
			"phy_wire_vlan_1":"` + ap.PhyWireVlan1 + `",
			"phy_wire_vlan_2":"` + ap.PhyWireVlan2 + `",
			"offline_selfmanage":"` + ap.OfflineSelfmanage + `",
			"username":"` + ap.Username + `",
			"password":"` + ap.Password + `",
			"led":"` + ap.Led + `",
			"led_wifi_sync":"` + ap.LedWifiSync + `",
			"led_timer_enable":"` + ap.LedTimerEnable + `"}
	}}`
	_, err := s.ds(ctx, body)
	return err
}
