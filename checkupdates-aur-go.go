package main

import (
	"fmt"
	"github.com/gobwas/glob"
)

func main() {
	packages := getForeignPackages()
	ignored := readConf("/etc/pacman.conf")
	getAurVersions(packages)
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
