// Copyright (c) 2019 Duxiaoman, Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller

import (
	"encoding/json"
	"github.com/chailyuan/lightsocks"
	"github.com/chailyuan/lightsocks/cmd"
	"github.com/chailyuan/lightsocks/server"
	"html/template"
	"log"
	"net/http"
	"sync"
)

//api接口固定格式
type Result struct {
	Ret    int
	Reason string
	Data   interface{}
}

//单例
type ApiController struct {
}

var once sync.Once
var apiController *ApiController

func GetApiController() *ApiController {
	once.Do(func() { apiController = &ApiController{} })
	return apiController
}

//controller路由和参数获取
func (this *ApiController) IndexAction(w http.ResponseWriter, r *http.Request, server *server.LsServer, config *cmd.Config) {
	defer r.Body.Close()
	w.Header().Set("content-type", "application/json")

	//获取所有任务列表
	if r.URL.Path == "/api/requestPass" {

		OutputJson(w, 0, "ok", config.Password)
		return
	}
	if r.URL.Path == "/api/changePass" {
		newPass := lightsocks.RandPassword()
		bsPassword, err := lightsocks.ParsePassword(newPass)
		if err != nil {
			OutputJson(w, 0, "ok", "解析密码异常")
		}

		cipher := lightsocks.GetInstance()
		cipher.SetPassword(bsPassword)

		config.Password = newPass
		config.SaveConfig()

		OutputJson(w, 0, "ok", newPass)
		return
	}

	//未路由到，跳404
	t, err := template.ParseFiles("static/404.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
}

//将结果打包json后返回
func OutputJson(w http.ResponseWriter, ret int, reason string, data interface{}) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")

	out := &Result{ret, reason, data}
	b, err := json.Marshal(out)
	if err != nil {
		return
	}
	//log.Println("finish api request result:", string(b))
	w.Write(b)
}
