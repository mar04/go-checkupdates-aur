package main

import (
	"fmt"
)

func main() {
	packages, err := getForeignPackages()
	if err != nil {
		return
	}
	getAurVersions(packages)
	printAurVersions(packages)
}

func printAurVersions(packages map[string]*pkgInfo) {
	for pkgName, v := range packages {
		if v.aurVersion != "" && isNewer(v.version, v.aurVersion) {
			fmt.Println(pkgName, v.version, "->", v.aurVersion)
		}
	}
}

/*
TODO: read /etc/pacman.conf and get IgnoredPackages
TODO: include ignored

*/
