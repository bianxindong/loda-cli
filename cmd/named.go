package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/lodastack/loda-cli/setting"
	"github.com/lodastack/models"
	"github.com/lodastack/sdk-go"
	"github.com/oiooj/cli"
)

var CmdNamed = cli.Command{
	Name:        "named",
	Usage:       "named config",
	Description: "named config",
	Action:      runNamed,
}

func runNamed(c *cli.Context) {
	header := `$ORIGIN .
$TTL 600	; 10 minutes
loda	IN SOA	ns1.loda ns.ifeng.com. (
				2006010215 ; serial
				600        ; refresh (10 minutes)
				600        ; retry (10 minutes)
				86400      ; expire (1 day)
				60         ; minimum (1 minute)
				)
			NS	ns1.loda.
			NS	ns2.loda.
$ORIGIN loda.
$TTL 60	; 1 minutes
ns1			A	10.80.40.157
ns2			A	10.90.1.225

`
	os.Remove("./loda.zone")
	var body string
	var nsList NameSpaceList
	ms, err := nsList.AllNameSpaces()
	if err != nil {
		handlerErr(err)
	}
	for _, ns := range ms {
		var serverList NamedServerList
		for _, server := range serverList.getServerList(ns, "machine") {
			body = fmt.Sprintf("%s%s		A 	%s\n", body, strings.TrimSuffix(ns, ".loda"), server.IP)
		}
	}
	if len(body) < 92484 {
		handlerErr(fmt.Errorf("%s", "body broken"))
	}
	f, err := os.Create("./loda.zone")
	if err != nil {
		handlerErr(err)
	}
	w := bufio.NewWriter(f)
	_, err = w.WriteString(header + body)
	if err != nil {
		handlerErr(err)
	}
	err = w.Flush()
	if err != nil {
		handlerErr(err)
	}
	Send(1)
}

func handlerErr(err error) {
	Send(0)
	fmt.Println(err)
	os.Exit(1)
}

type NamedServerList struct {
	Members   []NamedServer `json:"data"`
	NameSpace string
}

type NamedServer struct {
	Hostname   string `json:"hostname"`
	IP         string `json:"ip"`
	LastReport string `json:"lastReport"`
	Status     string `json:"status"`
	Version    string `json:"version"`
}

func (this *NamedServerList) getServerList(ns, resType string) []NamedServer {
	url := fmt.Sprintf(setting.API_Res, ns, resType)
	resp, err := http.Get(url)
	if err != nil {
		handlerErr(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		handlerErr(err)
	}
	json.Unmarshal(body, &this)
	if len(this.Members) == 0 {
		return this.Members
	}
	m := make(map[string]struct{})
	var res []NamedServer
	for _, s := range this.Members {
		for _, ip := range strings.Split(s.IP, ",") {
			if _, ok := m[ip]; ok {
				continue
			}
			if namedIsIntranet(ip) {
				s.IP = ip
				res = append(res, s)
				m[ip] = struct{}{}
				break
			}
		}
	}
	return res
}

func namedIsIntranet(ipStr string) bool {
	if strings.TrimSpace(ipStr) == "" || strings.Contains(ipStr, ",") {
		return false
	}
	if strings.HasPrefix(ipStr, "10.") {
		return true
	}
	return false
}

func Send(value float64) error {
	m := models.Metric{
		Name:      "named.sync",
		Timestamp: time.Now().Unix(),
		Value:     value,
	}
	data, err := json.Marshal([]models.Metric{m})
	if err != nil {
		return err
	}
	return sdk.Post("registry.monitor.loda", data)
}
