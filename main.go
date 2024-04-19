package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
	"tpset/assets"

	_ "time/tzdata"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

type Setting struct {
	APMac        string
	APName       string
	WiredVlan    string
	WirelessVlan string
	SSID         string
	Password     string
}

func init() {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	time.Local = loc
}

func main() {
	var (
		addr string
	)

	flag.StringVar(&addr, "l", "127.0.0.1:10000", "listen address")
	flag.Parse()

	debug := os.Getenv("DEBUG") == "1"

	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	{
		var fs static.ServeFileSystem
		if debug {
			fs = static.LocalFile("assets/html", false)
		} else {
			fs = static.EmbedFolder(assets.Html, "html")
		}
		r.Use(static.Serve("/", fs))
	}

	r.POST("/handle", handler())

	s := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	if runtime.GOOS == "darwin" {
		go func() {
			port := strings.LastIndex(addr, ":")
			if port < 0 {
				return
			}
			time.Sleep(time.Second)
			exec.Command("open", "http://localhost:"+addr[port+1:]).Run()
		}()
	}

	s.ListenAndServe()
}

func handler() gin.HandlerFunc {
	var lock sync.Mutex

	return func(ctx *gin.Context) {
		lock.Lock()
		defer lock.Unlock()

		ctx.Header("Transfer-Encoding", "chunked")
		logger := log.New(ctx.Writer, "", log.LstdFlags)

		ac := ctx.PostForm("ac")
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")
		file, err := ctx.FormFile("file")
		if err != nil {
			logger.Println(err)
			return
		}

		f, err := file.Open()
		if err != nil {
			logger.Println(err)
			return
		}
		defer f.Close()

		log.Println("processing csv file...")

		bufReader := bufio.NewReader(f)
		loop := true
		for i := 1; loop; i++ {
			line, err := bufReader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					loop = false
				} else {
					ctx.Writer.WriteString(err.Error())
					return
				}
			}

			if i == 1 {
				log.Println("skip first line as header")
				continue
			}

		}

		for range 30 {
			time.Sleep(time.Second)
			ctx.Writer.WriteString(time.Now().Format(time.DateTime) + "\n")
			ctx.Writer.Flush()
		}
	}
}

/*
func main1() {
	var (
		csv      string
		ac       string
		username string
		password string
	)

	flag.StringVar(&csv, "csv", "", "csv files")
	flag.StringVar(&ac, "ac", "", "ac address like http://192.168.1.1")
	flag.StringVar(&username, "username", "", "ac username")
	flag.StringVar(&password, "password", "", "ac password (encrypted), get from chrome network")
	flag.Parse()

	if csv == "" || ac == "" || username == "" || password == "" {
		flag.Usage()
		os.Exit(1)
	}

	data, err := os.ReadFile(csv)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(data), "\n")
	list := make([]Setting, 0)
	for i, line := range lines {
		if i == 0 {
			continue
		}
		items := strings.Split(line, ",")
		if len(items) != 6 {
			continue
		}
		var s Setting
		s.APMac = strings.TrimSpace(items[0])
		s.APName = strings.TrimSpace(items[1])
		s.WiredVlan = strings.TrimSpace(items[2])
		s.WirelessVlan = strings.TrimSpace(items[3])
		s.SSID = strings.TrimSpace(items[4])
		s.Password = strings.TrimSpace(items[5])
		s.APMac = tplink.SimplifyMAC(s.APMac)
		list = append(list, s)
	}

	s := tplink.NewService(ac)
	ctx := context.Background()
	timeout := 3 * time.Second

	func() {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		err := s.Login(ctx, username, password)
		if err != nil {
			log.Fatal(err)
		}
	}()

	var groupList []tplink.Group
	func() {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var err error
		list, err := s.ListGroup(ctx)
		if err != nil {
			log.Fatal(err)
		}
		groupList = list
	}()

	var apList []tplink.AP
	for _, group := range groupList {
		func() {
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			list, err := s.ListAP(ctx, group.GroupId)
			if err != nil {
				log.Fatal(err)
			}
			apList = append(apList, list...)
		}()
	}

	var ssidList []tplink.SSID
	func() {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		list, err := s.ListSSID(ctx)
		if err != nil {
			log.Fatal(err)
		}
		ssidList = list
	}()

	apmacs := bmap.NewMapFromSlice(apList, func(ap tplink.AP) string { return tplink.SimplifyMAC(ap.Mac) })
	ssids := bmap.NewMapFromSlice(ssidList, func(ssid tplink.SSID) string { return ssid.SSID })
	for _, set := range list {
		ap, ok := apmacs[set.APMac]
		if !ok {
			fmt.Printf("apmac: %v not found on ac\n", set.APMac)
			os.Exit(1)
		}

		ap.EntryName = set.APName
		ap.PhyWireVlan1 = set.WiredVlan
		ap.PhyWireVlan2 = set.WiredVlan
		func() {
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			err := s.SetAP(ctx, ap)
			if err != nil {
				log.Fatal(err)
			}
		}()

		ssidHandler := s.SetSSID
		ssid, ok := ssids[set.SSID]
		if !ok {
			ssidHandler = s.AddSSID
			ssid = tplink.SSID{}
		}

		ssid.Auth = "3"
		ssid.AutoBind = "off"
		ssid.BwCtrlEnable = "0"
		ssid.BwCtrlMode = "1"
		ssid.Cipher = "2"
		ssid.DefaultBindFreq = "771"
		ssid.DefaultBindVlan = "0"
		ssid.Desc = ""
		ssid.DownLimit = "128"
		ssid.Enable = "on"
		ssid.Encryption = "1"
		ssid.Isolate = "0"
		ssid.Key = set.Password
		ssid.KeyUpdateIntv = "86400"
		ssid.SSID = set.SSID
		ssid.SSIDCodeType = "1"
		ssid.SSIDBrd = "1"
		ssid.UpLimit = "128"

		func() {
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			err := ssidHandler(ctx, ssid)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	// refresh ap list
	apList = apList[:0]
	for _, group := range groupList {
		func() {
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			list, err := s.ListAP(ctx, group.GroupId)
			if err != nil {
				log.Fatal(err)
			}
			apList = append(apList, list...)
		}()
	}

	// refresh ssid list
	ssidList = ssidList[:0]
	func() {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		list, err := s.ListSSID(ctx)
		if err != nil {
			log.Fatal(err)
		}
		ssidList = list
	}()

	// get radioList
	var radioList []tplink.Radio
	func() {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		list, err := s.ListRadio(ctx)
		if err != nil {
			log.Fatal(err)
		}
		radioList = list
	}()

	ssids = bmap.NewMapFromSlice(ssidList, func(ssid tplink.SSID) string { return ssid.SSID })
	apnames := bmap.NewMapFromSlice(apList, func(ap tplink.AP) string { return ap.EntryName })
	apradios := make(map[string][]tplink.Radio)
	for _, radio := range radioList {
		mac := tplink.SimplifyMAC(apnames[radio.ApName].Mac)
		list, ok := apradios[mac]
		if !ok {
			list = make([]tplink.Radio, 0)
		}
		list = append(list, radio)
		apradios[mac] = list
	}

	for _, set := range list {
		radios, ok := apradios[set.APMac]
		if !ok {
			fmt.Printf("apmac %v no radio found", set.APMac)
			os.Exit(1)
		}

		ssid, ok := ssids[set.SSID]
		if !ok {
			fmt.Printf("ssid %v missing on ac", set.SSID)
			os.Exit(1)
		}

		radioIds := bmap.NewMapFromSlice(radios, func(radio tplink.Radio) string { return radio.RfId }).Keys()

		func() {
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			err := s.Bind(ctx, tplink.Bind{
				ServId:  ssid.ServId,
				VlanId:  set.WirelessVlan,
				RadioId: radioIds,
			})
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
	fmt.Println("all done.")
}
*/
