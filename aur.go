package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func constructQuery(packages map[string]*pkgInfo) string {
	const (
		aur     = "https://aur.archlinux.org"
		rpc     = "/rpc/?v=5"
		info    = "&type=info"
		infoArg = "&arg[]="
	)

	query := strings.Builder{}
	query.WriteString(aur + rpc + info)
	for pkgName := range packages {
		query.WriteString(infoArg + pkgName)
	}

	return query.String()
}

func getAurVersions(packages map[string]*pkgInfo) {
	query := constructQuery(packages)

	resp, err := http.Get(query)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	var r struct {
		Type    string
		Results []struct {
			Name    string
			Version string
		}
	}
	err = json.Unmarshal(respBody, &r)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	if r.Type != "multiinfo" {
		os.Exit(4)
	}

	for _, p := range r.Results {
		packages[p.Name].aurVersion = p.Version
	}
}
