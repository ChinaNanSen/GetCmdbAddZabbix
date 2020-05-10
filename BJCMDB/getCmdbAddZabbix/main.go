package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
	"strings"
	"sync"

	coll "github.com/chenhg5/collection"
	zabbixMt "github.com/zabbixoperating"
)

//定义主机服务信息
type ServiceInfo struct {
	DS string `json:"DS"`
}

//定义主机信息结构体
type HostInfo struct {
	Service  []ServiceInfo `json:"service"`
	Hostname string        `json:"hostname"`
	Vlanip   string        `json:"vlanip"`
}

type CmdbData struct {
	Objects []HostInfo `json:"objects"`
}

func httpPostForm() (info []byte) {
	//url := "http://root:babytree.com@portalapi.babytree.com/api/v1/allportal/?format=json&limit=5000"
	url := "http://root:babytree.com@10.1.7.11/api/v1/allportal/?format=json&limit=5000"
	resp, err := http.Get(url) //改送HTTP Get请求
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("状态码是：%v\n", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return body
}
func AddData(ch1 chan HostInfo) {
	var msg CmdbData
	err := json.Unmarshal(httpPostForm(), &msg)
	if err != nil {
		fmt.Printf("cmdb数据反序列化失败,%v\n", err)
	}
	fmt.Println("过滤不需要的数据，然后循环往通道插入数据")
	//最后关闭通道
	defer close(ch1)
	for _, v := range msg.Objects {
		//过滤不需要添加监控的机器
		if strings.Contains(v.Hostname, "onl1") == false && strings.Contains(v.Hostname, "onl2") == false {
			continue
		}
		ch1 <- v
	}
	wg.Done()
}

func UpdateData(ch1 chan HostInfo) {
	zabbixHost := zabbixMt.GetHost()
	fmt.Println("开始循环从通道中取数据,并阻塞通道")
	for v := range ch1 {
		if len(v.Service) == 0 {
			fmt.Println("空服务主机：", v.Hostname)
			continue
		}
		if zabbixMt.GetHostGroupID(v.Service[0].DS) == "" || zabbixMt.GetHostGroupID(v.Service[0].DS) == nil {
			_, err := zabbixMt.CreateHostGroup(v.Service[0].DS)
			if err != nil {
				fmt.Println(v.Service[0].DS, "服务组创建失败!")
			}
			fmt.Println(v.Service[0].DS, "服务组创建成功")
		}
		if coll.Collect(zabbixHost).Contains(v.Hostname) == false {
			relt, err := zabbixMt.CreateHost(v.Vlanip, v.Hostname, zabbixMt.GetHostGroupID(v.Service[0].DS).(string), "10001")
			if err != nil {
				fmt.Printf("创建主机失败-->%s", err)
				return
			}
			fmt.Println(v.Hostname, "创建主机成功", relt)
		}
	}
	wg.Done()
}

var wg sync.WaitGroup

func main() {
	f, err := os.Create("D:/Logs/cpu.pprof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	ch1 := make(chan HostInfo)
	wg.Add(2)
	go UpdateData(ch1)
	go AddData(ch1)
	wg.Wait()
}
