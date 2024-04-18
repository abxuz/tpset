package mercury

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
)

type Radio struct {
	Airtime              int    `json:"airtime"`
	Antenna              int    `json:"antenna"`
	ApName               string `json:"ap_name"`
	BcnIntvl             int    `json:"bcnIntvl"`
	BrdProbe             int    `json:"brd_probe"`
	BwMode               string `json:"bw_mode"`
	DftTxPwr             int    `json:"dftTxPwr"`
	DtimPrd              int    `json:"dtim_prd"`
	FragThr              int    `json:"frag_thr"`
	FreqName             string `json:"freq_name"`
	FreqBand             int    `json:"freqBand"`
	IsHdap               int    `json:"ishdap"`
	MaxTxPwr             int    `json:"maxTxPwr"`
	MaxUsr               int    `json:"maxUsr"`
	MinTxPwr             int    `json:"minTxPwr"`
	MutexChannel         string `json:"mutex_channel"`
	PwrLevel             int    `json:"pwr_level"`
	PwrLevelUi           int    `json:"pwr_level_ui"`
	RadioId              string `json:"radioId"`
	RdName               string `json:"rd_name"`
	RdMode               string `json:"rd_mode"`
	RdChnl               int    `json:"rd_chnl"`
	RdChnlDcs            int    `json:"rd_chnl_dcs"`
	RdChnlTolerance      int    `json:"rd_chnl_tolerance"`
	RdChnlOccupThreshold int    `json:"rd_chnl_occup_threshold"`
	RdCheckCircle        int    `json:"rd_check_circle"`
	RdOnlineSwitch       int    `json:"rd_online_switch"`
	RdState              string `json:"rd_state"`
	RdMaxUsr             int    `json:"rdMaxUsr"`
	RejectDiff           int    `json:"rejectDiff"`
	RejectEnable         int    `json:"rejectEnable"`
	RejectLimit          int    `json:"rejectLimit"`
	RssiRestrict         int    `json:"rssi_restrict"`
	RssiRestrictVal      int    `json:"rssiRestrictVal"`
	RssiKick             int    `json:"rssi_kick"`
	RssiKickVal          int    `json:"rssiKickVal"`
	RtsThr               int    `json:"rts_thr"`
	ShortGi              int    `json:"short_gi"`
	Supp11ac             int    `json:"supp11Ac"`
	Supp11ax             int    `json:"supp11Ax"`
	Supp160M             int    `json:"supp160m"`
	Wmm                  int    `json:"wmm"`
}

func (s *Service) ListRadio(ctx context.Context) ([]Radio, error) {
	uri := "/admin/ac_rdmngr?form=rdentry"
	_, err := s.request(ctx, uri, `{"method":"change","params":{"pageSize":"1000"}}`)
	if err != nil {
		return nil, err
	}

	list := make([]Radio, 0)
	for i := 0; true; i++ {
		data, err := s.request(ctx, uri, `{"method":"get","params":{"pageNo":`+strconv.Itoa(i)+`,"group_id":"-1"}}`)
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

		var arr []Radio
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
