package zabbixoperating

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/auth"
	"io/ioutil"
	"net/http"

)

var zabbixAuth = auth.ZabbixAuth()

//上海zabbix 3.0.5 API
//const URL string = "http://10.54.11.79/zabbix/api_jsonrpc.php"

//北京zabbix 4.4.7 API
const URL string = "http://172.26.10.197/api_jsonrpc.php"

type RequestBodys struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
	Auth    string      `json:"auth"`
}

type Response struct {
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	ID      int         `json:"id"`
}

var ResultData = Response{}
var err error

func NewRequestMethod(method string, params interface{}) interface{} {
	jsonObj := &RequestBodys{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
		Auth:    zabbixAuth,
	}
	encoded, err := json.Marshal(jsonObj)
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
	//最后关闭请求
	defer response.Body.Close()
	//读取请求的内容
	body, _ := ioutil.ReadAll(response.Body)
	//反序列化，将请求的内容初始化结构体
	err = json.Unmarshal(body, &ResultData)
	if err != nil {
		fmt.Printf("反序列化失败-->:%s", err)
	}
	return ResultData.Result
}

//创建主机组
func CreateHostGroup(grName string) (interface{}, error) {
	type Info struct {
		Name string `json:"name"`
	}
	data := &Info{}
	data.Name = grName
	Res := NewRequestMethod("hostgroup.create", data)
	return Res, err
}

//获取主机组ID
func GetHostGroupID(grName string) (ResultData interface{}) {
	type Fler struct {
		Name []string `json:"name"`
	}
	type Info struct {
		Output string `json:"output"`
		Filter Fler   `json:"filter"`
	}

	data := &Info{}
	data.Output = "extend"
	data.Filter = Fler{Name: []string{grName}}
	//fmt.Println(data)
	Res := NewRequestMethod("hostgroup.get", data)
	for _, v := range Res.([]interface{}) {
		s := v.(map[string]interface{})
		return s["groupid"]
	}
	return
}

//创建主机
func CreateHost(ip, hostname, groupid, tempid string) (interface{}, error) {
	type H struct {
		Type  int    `json:"type"`
		Main  int    `json:"main"`
		Useip int    `json:"useip"`
		IP    string `json:"ip"`
		Dns   string `json:"dns"`
		Port  string `json:"port"`
	}
	type Gid struct {
		Groupid string `json:"groupid"`
	}
	type Tid struct {
		Templateid string `json:"templateid"`
	}

	type Info struct {
		Host       string `json:"host"`
		Interfaces []H    `json:"interfaces"`
		Groups     []Gid  `json:"groups"`
		Templates  []Tid  `json:"templates"`
	}

	HostInfo := []H{
		{Type: 1, Main: 1, Useip: 1, IP: ip, Dns: "", Port: "10050"},
		//{Type: 1, Main: 1, Useip: 0, IP: hostname, Dns: hostname, Port: "10050"},
	}
	GroupInfo := []Gid{
		{Groupid: groupid},
	}
	TempInfo := []Tid{
		{Templateid: tempid},
	}
	data := &Info{}
	data.Host = hostname
	data.Interfaces = HostInfo
	data.Groups = GroupInfo
	data.Templates = TempInfo
	Res := NewRequestMethod("host.create", *data)
	return Res, err
}

//创建docker主机
func CreateDockerHost(hostname, groupid, tempid string) (interface{}, error) {
	type H struct {
		Type  int    `json:"type"`
		Main  int    `json:"main"`
		Useip int    `json:"useip"`
		IP    string `json:"ip"`
		Dns   string `json:"dns"`
		Port  string `json:"port"`
	}
	type Gid struct {
		Groupid string `json:"groupid"`
	}
	type Tid struct {
		Templateid string `json:"templateid"`
	}

	type Info struct {
		Host       string `json:"host"`
		Interfaces []H    `json:"interfaces"`
		Groups     []Gid  `json:"groups"`
		Templates  []Tid  `json:"templates"`
	}

	HostInfo := []H{
		//{Type: 1, Main: 1, Useip: 1, IP: ip, Dns: "", Port: "10050"},
		{Type: 1, Main: 1, Useip: 0, IP: "", Dns: hostname, Port: "10050"},
	}
	GroupInfo := []Gid{
		{Groupid: groupid},
	}
	TempInfo := []Tid{
		{Templateid: tempid},
	}
	data := &Info{}
	data.Host = hostname
	data.Interfaces = HostInfo
	data.Groups = GroupInfo
	data.Templates = TempInfo
	Res := NewRequestMethod("host.create", *data)
	return Res, err
}

//获得主机名
func GetHost() (ResultData interface{}) {
	type Info struct {
		Output []string `json:"output"`
	}
	data := Info{}
	data.Output = []string{"host"}
	Res := NewRequestMethod("host.get", &data)
	//fmt.Println(Res)
	var zabbixH []string
	for _, v := range Res.([]interface{}) {
		s := v.(map[string]interface{})
		zabbixH = append(zabbixH, s["host"].(string))
	}
	return zabbixH
}
