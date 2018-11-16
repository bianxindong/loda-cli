package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/oiooj/loda-cli/setting"
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

func UMachineInit() []string {
	res := MachineInit()
	if len(res) != 0 {
		for k, v := range res {
			res[k] = reverse(v)
		}
		return res
	} else {
		return []string{"error"}
	}
}

// reverse string by '.'
func reverse(s string) string {
	splitStr := strings.Split(s, ".")
	if len(splitStr) != 0 {
		for i, j := 0, len(splitStr)-1; i < j; i, j = i+1, j-1 {
			splitStr[i], splitStr[j] = splitStr[j], splitStr[i]
		}
		return strings.Replace(strings.Trim(fmt.Sprint(splitStr), "[]"), " ", ".", -1)
	} else {
		return "error"
	}
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

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func GetParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

//PathExist checks the pathfile is or isn't exist
func PathExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
