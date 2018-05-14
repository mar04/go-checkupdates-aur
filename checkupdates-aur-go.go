package main

import (
	"fmt"

	"github.com/gobwas/glob"
)

func main() {
	var ignored []string
	var ready = make(chan int)
	go getIgnored("/etc/pacman.conf", &ignored, ready)
	packages := getForeignPackages()
	getAurVersions(packages)
	<-ready
	printAurVersions(packages, ignored)
}

func printAurVersions(packages map[string]*pkgInfo, ignored []string) {
	for pkgName, v := range packages {
		if v.aurVersion != "" && isNewer(v.version, v.aurVersion) && notIgnored(pkgName, ignored) {
			fmt.Println(pkgName, v.version, "->", v.aurVersion)
		}
	}
}

func notIgnored(pkgName string, ignored []string) bool {
	for _, i := range ignored {
		g := glob.MustCompile(i)
		if g.Match(pkgName) {
			return false
		}
	}
	return true
}
