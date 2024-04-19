package main

import (
	"context"
	"strings"
	"time"
	"tpset/mercury"
	"tpset/tplink"
)

type AP interface {
	GetAPMac() string
	GetAPName() string
	SetAPName(name string)
	SetWiredVlan(vlan string)
	Clone() AP
}

type SSID interface {
	GetId() string
	GetSSID() string
	GetPassword() string
	SetPassword(pwd string)
	Clone() SSID
}

type Radio interface {
	GetId() string
	GetAPName() string
}

type Bind struct {
	ServId  string
	VlanId  string
	RadioId []string
}

type ACService interface {
	IsAC(ctx context.Context) (bool, error)
	Login(ctx context.Context, username, password string) error

	ListAP(ctx context.Context) ([]AP, error)
	SetAP(ctx context.Context, prev, cur AP) error

	ListSSID(ctx context.Context) ([]SSID, error)
	AddSSID(ctx context.Context, ssid SSID) error
	SetSSID(ctx context.Context, prev, cur SSID) error

	ListRadio(ctx context.Context) ([]Radio, error)
	Bind(ctx context.Context, bind Bind) error
}

func defaultTimeoutCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, time.Second*3)
}

type TPService struct {
	*tplink.Service
}

func (s *TPService) IsAC(ctx context.Context) (bool, error) {
	ctx, done := defaultTimeoutCtx(ctx)
	defer done()
	return s.Service.IsAC(ctx)
}

func (s *TPService) Login(ctx context.Context, username, password string) error {
	ctx, done := defaultTimeoutCtx(ctx)
	defer done()
	return s.Service.Login(ctx, username, password)
}

type TPAP struct {
	tplink.AP
}

func (ap *TPAP) GetAPMac() string {
	return ap.AP.Mac
}
func (ap *TPAP) GetAPName() string {
	return ap.AP.EntryName
}
func (ap *TPAP) SetAPName(name string) {
	ap.AP.EntryName = name
}
func (ap *TPAP) SetWiredVlan(vlan string) {
	ap.AP.PhyWireVlan1 = vlan
	ap.AP.PhyWireVlan2 = vlan
}
func (ap *TPAP) Clone() AP {
	return &TPAP{AP: ap.AP}
}

func (s *TPService) ListAP(ctx context.Context) ([]AP, error) {
	groups, err := s.Service.ListGroup(ctx)
	if err != nil {
		return nil, err
	}

	list := make([]AP, 0)
	for _, g := range groups {
		aps, err := s.Service.ListAP(ctx, g.GroupId)
		if err != nil {
			return nil, err
		}
		for _, ap := range aps {
			list = append(list, &TPAP{AP: ap})
		}
	}
	return list, nil
}

func (s *TPService) SetAP(ctx context.Context, prev, cur AP) error {
	return s.Service.SetAP(ctx, cur.(*TPAP).AP)
}

type TPSSID struct {
	tplink.SSID
}

func (ssid *TPSSID) GetId() string {
	return ssid.SSID.ServId
}
func (ssid *TPSSID) GetSSID() string {
	return ssid.SSID.SSID
}
func (ssid *TPSSID) GetPassword() string {
	return ssid.SSID.Key
}
func (ssid *TPSSID) SetPassword(pwd string) {
	ssid.SSID.Key = pwd
}
func (ssid *TPSSID) Clone() SSID {
	return &TPSSID{SSID: ssid.SSID}
}

func (s *TPService) ListSSID(ctx context.Context) ([]SSID, error) {
	ssids, err := s.Service.ListSSID(ctx)
	if err != nil {
		return nil, err
	}
	list := make([]SSID, 0)
	for _, ssid := range ssids {
		list = append(list, &TPSSID{SSID: ssid})
	}
	return list, nil
}

func (s *TPService) AddSSID(ctx context.Context, ssid SSID) error {
	var tpssid tplink.SSID

	tpssid.Auth = "3"
	tpssid.AutoBind = "off"
	tpssid.BwCtrlEnable = "0"
	tpssid.BwCtrlMode = "1"
	tpssid.Cipher = "2"
	tpssid.DefaultBindFreq = "771"
	tpssid.DefaultBindVlan = "0"
	tpssid.Desc = ""
	tpssid.DownLimit = "128"
	tpssid.Enable = "on"
	tpssid.Encryption = "1"
	tpssid.Isolate = "0"
	tpssid.Key = ssid.GetPassword()
	tpssid.KeyUpdateIntv = "86400"
	tpssid.SSID = ssid.GetSSID()
	tpssid.SSIDCodeType = "1"
	tpssid.SSIDBrd = "1"
	tpssid.UpLimit = "128"

	return s.Service.AddSSID(ctx, tpssid)
}
func (s *TPService) SetSSID(ctx context.Context, prev, cur SSID) error {
	return s.Service.SetSSID(ctx, cur.(*TPSSID).SSID)
}

type TPRadio struct {
	tplink.Radio
}

func (radio *TPRadio) GetId() string {
	return radio.Radio.RfId
}
func (radio *TPRadio) GetAPName() string {
	return radio.Radio.ApName
}

func (s *TPService) ListRadio(ctx context.Context) ([]Radio, error) {
	radios, err := s.Service.ListRadio(ctx)
	if err != nil {
		return nil, err
	}
	list := make([]Radio, 0)
	for _, radio := range radios {
		list = append(list, &TPRadio{Radio: radio})
	}
	return list, nil
}

func (s *TPService) Bind(ctx context.Context, bind Bind) error {
	return s.Service.Bind(ctx, tplink.Bind{
		ServId:  bind.ServId,
		VlanId:  bind.VlanId,
		RadioId: bind.RadioId,
	})
}

type MercuryService struct {
	*mercury.Service
}

func (s *MercuryService) IsAC(ctx context.Context) (bool, error) {
	ctx, done := defaultTimeoutCtx(ctx)
	defer done()
	return s.Service.IsAC(ctx)
}

func (s *MercuryService) Login(ctx context.Context, username, password string) error {
	ctx, done := defaultTimeoutCtx(ctx)
	defer done()
	return s.Service.Login(ctx, username, password)
}

type MercuryAP struct {
	mercury.AP
}

func (ap *MercuryAP) GetAPMac() string {
	return ap.AP.Mac
}
func (ap *MercuryAP) GetAPName() string {
	return ap.AP.EntryName
}
func (ap *MercuryAP) SetAPName(name string) {
	ap.AP.EntryName = name
}
func (ap *MercuryAP) SetWiredVlan(vlan string) {
	ap.AP.PhyWireVlan1 = vlan
	ap.AP.PhyWireVlan2 = vlan
}
func (ap *MercuryAP) Clone() AP {
	return &MercuryAP{AP: ap.AP}
}

func (s *MercuryService) ListAP(ctx context.Context) ([]AP, error) {
	groups, err := s.Service.ListGroup(ctx)
	if err != nil {
		return nil, err
	}

	list := make([]AP, 0)
	for _, g := range groups {
		aps, err := s.Service.ListAP(ctx, g.Id)
		if err != nil {
			return nil, err
		}
		for _, ap := range aps {
			list = append(list, &MercuryAP{AP: ap})
		}
	}
	return list, nil
}

func (s *MercuryService) SetAP(ctx context.Context, prev, cur AP) error {
	return s.Service.SetAP(ctx, prev.(*MercuryAP).AP, cur.(*MercuryAP).AP)
}

type MercurySSID struct {
	mercury.SSID
}

func (ssid *MercurySSID) GetId() string {
	return ssid.SSID.Id
}
func (ssid *MercurySSID) GetSSID() string {
	return ssid.SSID.SSID
}
func (ssid *MercurySSID) GetPassword() string {
	return ssid.SSID.Psk
}
func (ssid *MercurySSID) SetPassword(pwd string) {
	ssid.SSID.Psk = pwd
}
func (ssid *MercurySSID) Clone() SSID {
	return &MercurySSID{SSID: ssid.SSID}
}

func (s *MercuryService) ListSSID(ctx context.Context) ([]SSID, error) {
	ssids, err := s.Service.ListSSID(ctx)
	if err != nil {
		return nil, err
	}
	list := make([]SSID, 0)
	for _, ssid := range ssids {
		list = append(list, &MercurySSID{SSID: ssid})
	}
	return list, nil
}
func (s *MercuryService) AddSSID(ctx context.Context, ssid SSID) error {
	var mssid mercury.SSID

	mssid.Enable = "1"
	mssid.SecType = "1"
	mssid.SSID = ssid.GetSSID()
	mssid.Psk = ssid.GetPassword()

	return s.Service.AddSSID(ctx, mssid)
}
func (s *MercuryService) SetSSID(ctx context.Context, prev, cur SSID) error {
	return s.Service.SetSSID(ctx, prev.(*MercurySSID).SSID, cur.(*MercurySSID).SSID)
}

type MercuryRadio struct {
	mercury.Radio
}

func (radio *MercuryRadio) GetId() string {
	return radio.Radio.RadioId
}
func (radio *MercuryRadio) GetAPName() string {
	if strings.HasSuffix(radio.Radio.ApName, "-p") {
		return strings.TrimSuffix(radio.Radio.ApName, "-p")
	}

	if strings.HasSuffix(radio.Radio.ApName, "-t") {
		return strings.TrimSuffix(radio.Radio.ApName, "-t")
	}

	return radio.Radio.ApName
}

func (s *MercuryService) ListRadio(ctx context.Context) ([]Radio, error) {
	radios, err := s.Service.ListRadio(ctx)
	if err != nil {
		return nil, err
	}
	list := make([]Radio, 0)
	for _, radio := range radios {
		list = append(list, &MercuryRadio{Radio: radio})
	}
	return list, nil
}

func (s *MercuryService) Bind(ctx context.Context, bind Bind) error {
	return s.Service.Bind(ctx, mercury.Bind{
		ServId:  bind.ServId,
		VlanId:  bind.VlanId,
		RadioId: bind.RadioId,
	})
}
