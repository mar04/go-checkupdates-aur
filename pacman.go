package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type pkgInfo struct {
	version    string
	aurVersion string
}

// GET PACKAGES
func getForeignPackages() (map[string]*pkgInfo, error) {
	cmd := exec.Command("/bin/pacman", "-Qm")
	out, err := cmd.Output()
	if err != nil {
		return nil, errors.New("can't read pacman's output")
	}

	packages := make(map[string]*pkgInfo)
	scanner := bufio.NewScanner(bytes.NewReader(out))
	scanner.Split(bufio.ScanWords)

	for i, p := 0, ""; scanner.Scan(); i++ {
		if i%2 == 0 {
			p = scanner.Text()
		} else {
			packages[p] = &pkgInfo{version: scanner.Text()}
		}
	}
	return packages, nil
}

func isNewer(old, new string) bool {
	cmd := exec.Command("/bin/vercmp", old, new)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return false //TODO: exit?
	}
	r, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		fmt.Println(err)
		return false //TODO:exit?
	}
	return r < 0
}

func getIgnored() map[string]int {
	ignored := make(map[string]int)

	//where is pacman.conf
	//open pacman.conf
	//follow includes?
	//get lines with ignoredpkg
	//take ignore groups into consideration?
	// tokenize ignored packages and apply shell glob? or store the package with globing characters?

	return ignored
}
