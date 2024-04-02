package tplink

import (
	"context"
	"encoding/json"
	"errors"
)

type Radio struct {
	Airtime              string `json:"airtime"`
	Antenna              string `json:"antenna"`
	ApName               string `json:"ap_name"`
	ApType               string `json:"ap_type"`
	Bandwidth            string `json:"bandwidth"`
	BcnIntvl             string `json:"bcn_intvl"`
	BrdProbe             string `json:"brd_probe"`
	Channel              string `json:"channel"`
	DftTxPwr             string `json:"dft_tx_pwr"`
	DtimPrd              string `json:"dtim_prd"`
	FragThr              string `json:"frag_thr"`
	FreqName             string `json:"freq_name"`
	FreqUnit             string `json:"freq_unit"`
	IsHdap               string `json:"is_hdap"`
	ManageFrameRate      string `json:"manage_frame_rate"`
	MaxTxPwr             string `json:"max_tx_pwr"`
	MaxUsr               string `json:"max_usr"`
	MinTxPwr             string `json:"min_tx_pwr"`
	Mode                 string `json:"mode"`
	MutexChannel         string `json:"mutex_channel"`
	Power                string `json:"power"`
	RejectDiff           string `json:"reject_diff"`
	RejectEnable         string `json:"reject_enable"`
	RejectLimit          string `json:"reject_limit"`
	RfCheckCircle        string `json:"rf_check_circle"`
	RfChnlDcs            string `json:"rf_chnl_dcs"`
	RfChnlOccupThreshold string `json:"rf_chnl_occup_threshold"`
	RfChnlTolerance      string `json:"rf_chnl_tolerance"`
	RfId                 string `json:"rf_id"`
	RfMaxUsr             string `json:"rf_max_usr"`
	RfName               string `json:"rf_name"`
	RfOnlineSwitch       string `json:"rf_online_switch"`
	RfState              string `json:"rf_state"`
	RssiRestrict         string `json:"rssi_restrict"`
	RssiRestrictVal      string `json:"rssi_restrict_val"`
	RtsThr               string `json:"rts_thr"`
	ShortGi              string `json:"short_gi"`
	Supp11ac             string `json:"supp_11ac"`
	Supp11ax             string `json:"supp_11ax"`
	Supp160M             string `json:"supp_160M"`
	Wmm                  string `json:"wmm"`
}

func (s *Service) ListRadio(ctx context.Context) ([]Radio, error) {
	body := `{
		"method":"get",
		"apmng_rf":{"table":"rf_entry","filter":[{"group_id":["-1"]}],
		"para":{"start":0,"end":1999}
	}}`
	data, err := s.ds(ctx, body)
	if err != nil {
		return nil, err
	}
	retErr := errors.New(string(data))

	var result struct {
		ApMngRf struct {
			RfEntry []map[string]json.RawMessage `json:"rf_entry"`
		} `json:"apmng_rf"`
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, retErr
	}

	list := make([]Radio, 0)
	for _, v := range result.ApMngRf.RfEntry {
		if len(v) != 1 {
			return nil, retErr
		}

		var raw json.RawMessage
		for _, raw = range v {
			break
		}

		var radio Radio
		err = json.Unmarshal(raw, &radio)
		if err != nil {
			return nil, retErr
		}
		list = append(list, radio)
	}

	return list, nil
}
