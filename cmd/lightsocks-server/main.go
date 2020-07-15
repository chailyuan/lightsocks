package main

import (
	"fmt"
	"github.com/chailyuan/lightsocks/controller"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/chailyuan/lightsocks"
	"github.com/chailyuan/lightsocks/cmd"
	"github.com/chailyuan/lightsocks/server"
	"github.com/phayes/freeport"
)

var version = "master"

var lsServer *server.LsServer
var config *cmd.Config

func main() {
	log.SetFlags(log.Lshortfile)

	// 优先从环境变量中获取监听端口
	port, err := strconv.Atoi(os.Getenv("LIGHTSOCKS_SERVER_PORT"))
	// 服务端监听端口随机生成
	if err != nil {
		port, err = freeport.GetFreePort()
	}
	if err != nil {
		// 随机端口失败就采用 7448
		port = 7448
	}
	// 默认配置
	config = &cmd.Config{
		ListenAddr: fmt.Sprintf(":%d", port),
		// 密码随机生成
		Password: lightsocks.RandPassword(),
	}
	config.ReadConfig()
	config.SaveConfig()

	go httpServer()

	// 启动 server 端并监听
	lsServer, err = server.NewLsServer(config.Password, config.ListenAddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(lsServer.Listen(func(listenAddr *net.TCPAddr) {
		log.Println(fmt.Sprintf(`
lightsocks-server:%s 启动成功，配置如下：
服务监听地址：
%s
密码：
%s`, version, listenAddr, config.Password))
	}))
}

func httpServer() {
	http.HandleFunc("/api/", IndexHandler) //设置访问的路由

	//启动http服务
	err := http.ListenAndServe(":12392", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//api路由进入controller
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("get api request:", r.URL.String())
	controller.GetApiController().IndexAction(w, r, lsServer, config)
}
