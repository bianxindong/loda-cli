package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"

	"github.com/lodastack/loda-cli/setting"
)

func MachineInit() []string {
	var nsList NameSpaceList
	res, err := nsList.AllNameSpaces()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return res
}

type NameSpaceList struct {
	Code    int      `json:"httpstatus"`
	Members []string `json:"data"`
}

func (this NameSpaceList) AllNameSpaces() ([]string, error) {
	resp, err := http.Get(setting.API_NS)
	if err != nil {
		return this.Members, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return this.Members, err
	}
	err = json.Unmarshal(body, &this)
	if err != nil {
		return this.Members, err
	}
	if len(this.Members) == 0 {
		fmt.Println("No NameSpace found!")
		return this.Members, nil
	}
	sort.Strings(this.Members)
	return this.Members, nil
}
