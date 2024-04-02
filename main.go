package main

import (
	"context"
	"log"
	"time"
	"tpset/tplink"
)

func main() {
	ctx := context.Background()
	s := tplink.NewService("http://100.64.63.42:6080")

	func() {
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		err := s.Login(ctx, "admin", "7rHLU7xB3vBzn4K")
		if err != nil {
			log.Fatal(err)
		}
	}()

	// var groupList []tplink.Group
	// func() {
	// 	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	// 	defer cancel()
	// 	list, err := s.ListGroup(ctx)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	groupList = list
	// }()

	// var apList []tplink.AP
	// for _, g := range groupList {
	// 	func() {
	// 		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	// 		defer cancel()
	// 		list, err := s.ListAP(ctx, g.GroupId)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		apList = append(apList, list...)
	// 	}()
	// }

	// for _, ap := range apList {
	// 	func() {
	// 		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	// 		defer cancel()
	// 		ap.EntryName += "y"
	// 		err := s.SetAP(ctx, ap)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 	}()
	// 	break
	// }

	var ssidList []tplink.SSID
	func() {
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		list, err := s.ListSSID(ctx)
		if err != nil {
			log.Fatal(err)
		}
		ssidList = list
	}()

	for _, ssid := range ssidList {
		if ssid.ServId != "47" {
			continue
		}

		func() {
			ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
			defer cancel()
			ssid.SSID += "x"
			err := s.SetSSID(ctx, ssid)
			if err != nil {
				log.Fatal(err)
			}
		}()
		break
	}

}
