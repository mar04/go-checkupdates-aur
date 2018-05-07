package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	packages, err := getForeignPackages()
	if err != nil {
		return
	}
	getAurVersions(packages)
	printAurVersions(packages)
}

func printAurVersions(packages map[string]pkgInfo) {
	for k, v := range packages {
		if v.aurVersion != "" && isNewer(v.version, v.aurVersion) {
			fmt.Println(k, v.version, "->", v.aurVersion)
		}
	}
}

func isNewer(old, new string) bool {
	cmd := exec.Command("/bin/vercmp", old, new)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return false
	}
	r, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		fmt.Println(err)
		return false
	}
	return r < 0
}

func constructQuery(packages map[string]pkgInfo) string {
	const (
		aur       = "https://aur.archlinux.org"
		rpc       = "/rpc/?v=5"
		search    = "&type=search"
		searchArg = "&arg="
		info      = "&type=info"
		infoArg   = "&arg[]="
	)

	query := strings.Builder{}
	query.WriteString(aur + rpc + info)
	for pkgName := range packages {
		query.WriteString(infoArg + pkgName)
	}

	return query.String()
}

func getAurVersions(packages map[string]pkgInfo) {
	query := constructQuery(packages)
	// fmt.Println(query)

	resp, err := http.Get(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var f interface{}
	err = json.Unmarshal(b, &f)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println("as json:", f)

	m := f.(map[string]interface{})
	if m["type"] != "multiinfo" {
		fmt.Printf("duba blada")
		return
	}
	z := m["results"].([]interface{})
	for _, v := range z {
		//fmt.Println(i,v)
		aurPkg := v.(map[string]interface{})
		aurName := aurPkg["Name"].(string)
		aurVer := aurPkg["Version"].(string)
		// fmt.Println(aurName, aurVer)
		packages[aurName] = pkgInfo{version: packages[aurName].version, aurVersion: aurVer} //TODO:ugly as hell
	}
}

/*
TODO: read /etc/pacman.conf and get IgnoredPackages
TODO: include ignored

*/
