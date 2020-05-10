package auth

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"fmt"
)

//package main
//
//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"io/ioutil"
//	"net/http"
//)

//将URL常量写入函数传参得方式
//上海zabbix 3.0.5 API
//const URL string = "http://10.54.11.79/zabbix/api_jsonrpc.php"

//北京zabbix 4.4.7 API
const URL string = "http://172.26.10.197/api_jsonrpc.php"

type UserInfo struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

var USER = UserInfo{
	User:     "zhangzhinan",
	Password: "zzn920617",
}

type RequestBodys struct {
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  UserInfo `json:"params"`
	ID      int      `json:"id"`
}

type Response struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"result"`
	ID      int    `json:"id"`
}

var Jsonobj = RequestBodys{
	Jsonrpc: "2.0",
	Method:  "user.login",
	Params:  USER,
	ID:      1,
}

func ZabbixAuth() (Result string) {
	//实例化
	resultData := Response{}
	//把结构体序列化 []byte 切片类型
	encoded, err := json.Marshal(Jsonobj)
	if err != nil {
		fmt.Println(err)
	}
	client := &http.Client{}
	//提交请求
	reqs, err := http.NewRequest("POST", URL, bytes.NewBuffer(encoded))
	//增加header选项
	reqs.Header.Add("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	//处理返回结果
	response, err := client.Do(reqs)
	if err != nil {
		fmt.Printf("处理返回失败", err)
	}
	//fmt.Println(response)
	//最后关闭请求
	defer response.Body.Close()
	//读取请求的内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("读取body失败,",err)
	}
	//fmt.Println(string(body))
	//反序列化，将请求的内容初始化结构体
	err = json.Unmarshal(body, &resultData)
	if err != nil {
		fmt.Println(err)
	}
	return resultData.Result
}

//func main()  {
//	a :=ZabbixAuth()
//	fmt.Println(a)
//}