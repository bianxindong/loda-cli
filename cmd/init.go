package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/lodastack/loda-cli/setting"
)

func MachineInit() []string {
	var nsList NameSpaceList
	return nsList.AllNameSpaces()
}

type NameSpaceList struct {
	Code    int      `json:"httpstatus"`
	Members []string `json:"data"`
}

func (this NameSpaceList) AllNameSpaces() []string {
	resp, err := http.Get(setting.API_NS)
	if err != nil {
		fmt.Println(err)
		return this.Members
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return this.Members
	}
	err = json.Unmarshal(body, &this)
	if err != nil {
		fmt.Println(err)
		return this.Members
	}
	if len(this.Members) == 0 {
		fmt.Println("No NameSpace found!")
		return this.Members
	}
	sort.Strings(this.Members)
	return this.Members
}
