package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	var ignored []string
	var ready = make(chan int)
	go getIgnored("/etc/pacman.conf", &ignored, ready)
	packages := getForeignPackages()
	getAurVersions(packages)
	<-ready
	printUpdates(packages, ignored)
}

func printUpdates(packages map[string]*pkgInfo, ignored []string) {
	for pkgName, v := range packages {
		if v.aurVersion != "" && isNewer(v.version, v.aurVersion) && !isIgnored(pkgName, ignored) {
			fmt.Println(pkgName, v.version, "->", v.aurVersion)
		}
	}
}

func isIgnored(pkgName string, ignored []string) bool {
	for _, pkgGlob := range ignored {
		if result, err := filepath.Match(pkgGlob, pkgName); err == nil && result {
			return true
		}
	}
	return false
}
