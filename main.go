package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
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
	"tpset/mercury"
	"tpset/tplink"

	_ "time/tzdata"

	"github.com/abxuz/b-tools/bmap"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

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

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		if err := s.ListenAndServe(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	if runtime.GOOS == "darwin" {
		go func() {
			port := strings.LastIndex(addr, ":")
			if port < 0 {
				return
			}

			timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), time.Second)
			defer timeoutCancel()

			select {
			case <-ctx.Done():
				return
			case <-timeoutCtx.Done():
			}
			exec.Command("open", "http://localhost:"+addr[port+1:]).Run()
		}()
	}

	<-ctx.Done()
}

type mySSID struct {
	ssid     string
	password string
}

func (s *mySSID) GetId() string {
	return ""
}

func (s *mySSID) GetSSID() string {
	return s.ssid
}

func (s *mySSID) GetPassword() string {
	return s.password
}

func (s *mySSID) SetPassword(pwd string) {
	s.password = pwd
}

func (s *mySSID) Clone() SSID {
	newSSID := *s
	return &newSSID
}

func handler() gin.HandlerFunc {
	var lock sync.Mutex

	type Setting struct {
		APMac        string
		APName       string
		WiredVlan    string
		WirelessVlan string
		SSID         string
		Password     string
	}

	return func(ctx *gin.Context) {
		lock.Lock()
		defer lock.Unlock()

		ctx.Header("Transfer-Encoding", "chunked")
		log := log.New(&flushWriter{ResponseWriter: ctx.Writer}, "", log.LstdFlags)

		ac := ctx.PostForm("ac")
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")
		file, err := ctx.FormFile("file")
		if err != nil {
			log.Println(err)
			return
		}

		f, err := file.Open()
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()

		log.Println("开始读取CSV文件...")

		bufReader := bufio.NewReader(f)
		loop := true
		settings := make([]Setting, 0)
		for i := 1; loop; i++ {
			line, err := bufReader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					loop = false
				} else {
					log.Println(err)
					return
				}
			}

			if i == 1 {
				log.Println("第1行标题行，跳过")
				continue
			}

			line = strings.TrimSpace(line)
			if line == "" {
				log.Printf("第%v行空行，跳过", i)
				continue
			}

			items := strings.Split(line, ",")
			if len(items) != 6 {
				log.Printf("第%v行格式错误，`%v`", i, line)
				return
			}

			var s Setting
			s.APMac = strings.TrimSpace(items[0])
			s.APName = strings.TrimSpace(items[1])
			s.WiredVlan = strings.TrimSpace(items[2])
			s.WirelessVlan = strings.TrimSpace(items[3])
			s.SSID = strings.TrimSpace(items[4])
			s.Password = strings.TrimSpace(items[5])
			s.APMac = SimplifyMAC(s.APMac)
			settings = append(settings, s)
		}

		log.Printf("读取到%v条记录", len(settings))

		log.Printf("开始判断AC类型...")

		var service ACService

		service = &TPService{Service: tplink.NewService(ac)}
		ok, err := service.IsAC(ctx)
		if err != nil {
			log.Println(err)
			return
		}

		if !ok {
			service = &MercuryService{Service: mercury.NewService(ac)}
			ok, err := service.IsAC(ctx)
			if err != nil {
				log.Println(err)
				return
			}
			if !ok {
				log.Println("未检测到TP-Link或水星的AC")
				return
			}
			log.Println("检测到AC类型为水星")
		} else {
			log.Println("检测到AC类型为TP-Link")
		}

		log.Println("开始登陆AC...")
		err = service.Login(ctx, username, password)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("登陆AC成功！！！请不要再对AC有任何操作了！！！")

		log.Println("开始获取AP列表...")
		apList, err := service.ListAP(ctx)
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("获取到%v个AP", len(apList))

		log.Println("开始设置APName/有线VLAN...")
		apmacs := bmap.NewMapFromSlice(apList, func(ap AP) string { return SimplifyMAC(ap.GetAPMac()) })
		for i, set := range settings {
			ap, ok := apmacs[set.APMac]
			if !ok {
				log.Printf("未在AC上找到APMAC [%v] 的AP，跳过 (%v/%v)", set.APMac, i+1, len(settings))
				continue
			}

			newAP := ap.Clone()
			newAP.SetAPName(set.APName)
			newAP.SetWiredVlan(set.WiredVlan)
			err = service.SetAP(ctx, ap, newAP)
			if err != nil {
				log.Println(err)
				return
			}
			log.Printf("[%v] 设置APName: %v, VlanId: %v，成功 (%v/%v)", set.APMac, set.APName, set.WiredVlan, i+1, len(settings))
		}

		log.Println("开始获取SSID列表...")
		ssidList, err := service.ListSSID(ctx)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("获取到%v个SSID", len(ssidList))

		log.Println("开始设置SSID...")
		ssids := bmap.NewMapFromSlice(ssidList, func(ssid SSID) string { return ssid.GetSSID() })
		for i, set := range settings {
			var act string
			ssid, ok := ssids[set.SSID]
			if !ok {
				ssid = &mySSID{
					ssid:     set.SSID,
					password: set.Password,
				}
				err = service.AddSSID(ctx, ssid)
				act = "新增"
			} else {
				newSSID := ssid.Clone()
				newSSID.SetPassword(set.Password)
				err = service.SetSSID(ctx, ssid, newSSID)
				act = "修改"
			}

			if err != nil {
				log.Println(err)
				return
			}
			log.Printf("%v SSID [%v] 密码: %v ，成功 (%v/%v)", act, set.SSID, set.Password, i+1, len(settings))
		}

		log.Println("重新获取AP列表...")
		apList, err = service.ListAP(ctx)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("重新获取SSID列表...")
		ssidList, err = service.ListSSID(ctx)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("获取射频列表...")
		radioList, err := service.ListRadio(ctx)
		if err != nil {
			log.Println(err)
			return
		}

		ssids = bmap.NewMapFromSlice(ssidList, func(ssid SSID) string { return ssid.GetSSID() })
		apnames := bmap.NewMapFromSlice(apList, func(ap AP) string { return ap.GetAPName() })
		apradios := make(map[string][]Radio)
		for _, radio := range radioList {
			apmac := SimplifyMAC(apnames[radio.GetAPName()].GetAPMac())
			list, ok := apradios[apmac]
			if !ok {
				list = make([]Radio, 0)
			}
			list = append(list, radio)
			apradios[apmac] = list
		}

		log.Println("开始将SSID绑定到AP的射频...")
		for i, set := range settings {
			radios, ok := apradios[set.APMac]
			if !ok {
				log.Printf("未找到APMAC [%v] 的射频，可能AP没上线，跳过 (%v/%v)", set.APMac, i+1, len(settings))
				continue
			}

			ssid, ok := ssids[set.SSID]
			if !ok {
				log.Printf("没找到SSID [%v]，创建没成功？", set.SSID)
				return
			}

			radioIds := bmap.NewMapFromSlice(radios, func(radio Radio) string { return radio.GetId() }).Keys()
			err = service.Bind(ctx, Bind{
				ServId:  ssid.GetId(),
				VlanId:  set.WirelessVlan,
				RadioId: radioIds,
			})
			if err != nil {
				log.Println(err)
				return
			}
			log.Printf("SSID [%v] 绑定到 APMAC [%v]，成功 (%v/%v)", set.SSID, set.APMac, i+1, len(settings))
		}

		log.Println("大功告成!")
	}
}

func SimplifyMAC(mac string) string {
	mac = strings.ToLower(mac)
	mac = strings.ReplaceAll(mac, "-", "")
	mac = strings.ReplaceAll(mac, ":", "")
	return mac
}

type flushWriter struct {
	gin.ResponseWriter
}

func (w *flushWriter) Write(b []byte) (int, error) {
	defer w.ResponseWriter.Flush()
	return w.ResponseWriter.Write(b)
}
