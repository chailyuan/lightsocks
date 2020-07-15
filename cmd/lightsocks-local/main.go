package main

import (
	"encoding/json"
	"fmt"
	"github.com/chailyuan/lightsocks"
	"github.com/chailyuan/lightsocks/controller"
	"log"
	"net"
	"time"

	"github.com/chailyuan/lightsocks/cmd"
	"github.com/chailyuan/lightsocks/local"
	"github.com/valyala/fasthttp"
)

var config *cmd.Config
var lsLocal *local.LsLocal

const (
	DefaultListenAddr = ":7448"
)

var version = "master"

func main() {
	log.SetFlags(log.Lshortfile)

	// 默认配置
	config = &cmd.Config{
		ListenAddr: DefaultListenAddr,
	}
	config.ReadConfig()
	config.SaveConfig()

	// 启动 local 端并监听
	var err error
	lsLocal, err = local.NewLsLocal(config.Password, config.ListenAddr, config.RemoteAddr)
	if err != nil {
		log.Fatalln(err)
	}

	go startTimer(getPass)

	log.Fatalln(lsLocal.Listen(func(listenAddr *net.TCPAddr) {
		log.Println(fmt.Sprintf(`
lightsocks-local:%s 启动成功，配置如下：
本地监听地址：
%s
远程服务地址：
%s
密码：
%s`, version, listenAddr, config.RemoteAddr, config.Password))
	}))
}

func startTimer(f func()) {
	go func() {
		for {
			f()
			now := time.Now()
			// 计算下一个零点
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
		}
	}()
}

func getPass() {
	proxyReq, proxyResp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(proxyReq)
	defer fasthttp.ReleaseResponse(proxyResp)

	// 默认是application/x-www-form-urlencoded
	proxyReq.Header.SetContentType("application/json")
	proxyReq.Header.SetMethod(fasthttp.MethodGet)

	proxyReq.SetRequestURI("http://96.43.93.116:12392/api/changePass")

	if err := fasthttp.Do(proxyReq, proxyResp); err != nil {
		log.Println("请求失败:", err.Error())
		return
	}

	result := controller.Result{}
	err := json.Unmarshal(proxyResp.Body(), &result)
	if err != nil {
		log.Println("json解析错误")
	}
	log.Println(result.Data)

	config.Password = fmt.Sprintf("%v", result.Data)
	config.SaveConfig()

	bsPassword, err := lightsocks.ParsePassword(config.Password)
	if err != nil {

	}

	cipher := lightsocks.NewCipher(bsPassword)

	lsLocal.Cipher = cipher
	config.ReadConfig()
}
