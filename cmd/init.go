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

// MachineInit init machine
func MachineInit() []string {
	var nsList NameSpaceList
	res, err := nsList.AllNameSpaces()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return res
}

// NameSpaceList struct
type NameSpaceList struct {
	Code    int      `json:"httpstatus"`
	Members []string `json:"data"`
}

// AllNameSpaces returns all ns
func (nl NameSpaceList) AllNameSpaces() ([]string, error) {
	resp, err := http.Get(setting.API_NS)
	if err != nil {
		return nl.Members, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nl.Members, err
	}
	err = json.Unmarshal(body, &nl)
	if err != nil {
		return nl.Members, err
	}
	if len(nl.Members) == 0 {
		fmt.Println("No NameSpace found!")
		return nl.Members, nil
	}
	sort.Strings(nl.Members)
	return nl.Members, nil
}
