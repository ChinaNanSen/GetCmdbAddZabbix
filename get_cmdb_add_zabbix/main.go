package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	coll "github.com/chenhg5/collection"
	zabbixMt "github.com/zabbixoperating"
)

//定义主机信息结构体
type CmdbInfo struct {
	ID          string
	Hostname    string
	IP          string
	SaltID      string
	Type        string
	Usage       string
	State       string
	Application string
	Stage       string
}

type CmdbData struct {
	Count    int
	Next     interface{}
	Previous interface{}
	Results  []CmdbInfo
}

func httpPostForm() (info []byte) {
	//url := "http://apollo.mtmm.com/apiv1/united_devices/?page_size=600"
	url := "http://apollo.mtmm.com/apiv1/united_devices/?page_size=50"
	resp, err := http.Get(url) //改送HTTP Get请求
	if err != nil {
		fmt.Println(err.Error())
		return
	}
    //fmt.Println(resp.Body)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("状态码是：%v", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//fmt.Println(body)
	return body
}
func main() {
	var msg CmdbData
	err := json.Unmarshal(httpPostForm(), &msg)
	if err != nil {
		fmt.Printf("反序列化失败,%v", err)
	}
	zabbixHost := zabbixMt.GetHost()
	//fmt.Println(zabbixHost)
	for _, v := range msg.Results {
		//过滤不需要添加监控的机器
		if strings.Contains(v.Hostname, "qa-vm") || strings.Contains(v.Hostname, "performance") ||
			strings.Contains(v.Hostname, "压测") || strings.Contains(v.Hostname, "onl-pingan") ||
			strings.Contains(v.Hostname, "h5-webpage-agent") {
			continue
		}
		//判断添加主机组和主机
		if strings.HasSuffix(v.Hostname, "alhz") && coll.Collect(zabbixHost).Contains(v.Hostname) == false {
			//fmt.Println(v.Hostname, v.IP, v.Application)
			_, err := zabbixMt.CreateHostGroup(v.Application)
			if err != nil {
				fmt.Println(v.Application, "服务组创建失败!")
				return
			}
			if v.Application == "" {
				host := strings.Split(v.Hostname, ".")
				v.Application = host[0]
				_, err := zabbixMt.CreateHostGroup(v.Application)
				if err != nil {
					fmt.Println(v.Application, "服务组创建再次失败", err)
				}
			}
			fmt.Println("创建服务组成功！！！", v.Application)
			Ggr := zabbixMt.GetHostGroupID(v.Application)
			values, ok := Ggr.(string)
			if !ok {
				fmt.Println("不是字符串类型")
				continue
			}
			//添加docker主机和主机组，模板不一样
			if v.Type == "DOCKER" {
				relts, errs := zabbixMt.CreateDockerHost(v.Hostname, values, "10106")
				if err != nil {
					fmt.Printf("创建主机失败-->%s", errs)
					return
				}
				fmt.Println("创建docker容器成功", relts)
			}
			relt, err := zabbixMt.CreateHost(v.IP, v.Hostname, values, "10001")
			if err != nil {
				fmt.Printf("创建主机失败-->%s", err)
				return
			}
			fmt.Println("创建主机成功", relt)
		}
	}
}
