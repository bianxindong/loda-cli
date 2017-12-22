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
	"github.com/oiooj/cli"
)

var CmdNamed = cli.Command{
	Name:        "named",
	Usage:       "named config",
	Description: "named config",
	Action:      runNamed,
}

func runNamed(c *cli.Context) {
	header := `$TTL	60 ; 24 hours could have been written as 24h or 1D
$ORIGIN loda.
; line below expands to: localhost 1D IN SOA localhost root.localhost
@  1D  IN	 SOA @	hostmaster (
					%s ; serial
					3H ; refresh
					15 ; retry
					1w ; expire
					3h ; minimum
					)


			NS		ns1
`
	t := time.Now()
	ts := t.Format("2006010215")
	header = fmt.Sprintf(header, ts)
	os.Remove("./loda.zone")
	var body string
	var nsList NameSpaceList
	for _, ns := range nsList.AllNameSpaces() {
		var serverList NamedServerList
		for _, server := range serverList.getServerList(ns, "machine") {
			body = fmt.Sprintf("%s%s		IN	A 	%s\n", body, strings.TrimSuffix(ns, ".loda"), server.IP)
		}
	}
	if len(body) < 92484 {
		fmt.Println("body broken")
		os.Exit(1)
	}
	f, err := os.Create("./loda.zone")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	w := bufio.NewWriter(f)
	_, err = w.WriteString(header + body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	w.Flush()
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
		fmt.Println("Get from loda error: ", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read from HTTP error: ", err)
		os.Exit(1)
	}
	json.Unmarshal(body, &this)
	if len(this.Members) == 0 {
		return this.Members
	}
	m := make(map[string]struct{})
	for i, s := range this.Members {
		for _, ip := range strings.Split(s.IP, ",") {
			if _, ok := m[ip]; ok {
				continue
			}
			if namedIsIntranet(ip) {
				s.IP = ip
				this.Members[i] = s
				m[ip] = struct{}{}
				break
			}
		}
	}
	return this.Members
}

func namedIsIntranet(ipStr string) bool {
	if strings.HasPrefix(ipStr, "10.") {
		return true
	}
	return false
}
