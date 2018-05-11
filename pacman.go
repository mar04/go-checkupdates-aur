package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type pkgInfo struct {
	version    string
	aurVersion string
}

// GET PACKAGES
func getForeignPackages() map[string]*pkgInfo {
	cmd := exec.Command("/bin/pacman", "-Qm")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(99)
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
	return packages
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

func readConf(pacmanConf string) []string {
	f, err := os.Open(pacmanConf)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	conf := bufio.NewReader(f)

	var ignored []string

	// fmt.Println("reading from file: ", pacmanConf)
	for {
		//read a line
		line, err := conf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "Include") {
			//strip comments
			if i := strings.IndexByte(line, '#'); i >= 0 {
				line = line[:i]
			}
			i := strings.IndexByte(line, '=')
			ignored = append(ignored, readConf(strings.Fields(line[i+1:])[0])...)
			continue
		}
		if strings.HasPrefix(line, "IgnorePkg") {
			//strip comments
			if i := strings.IndexByte(line, '#'); i >= 0 {
				line = line[:i]
			}
			i := strings.IndexByte(line, '=')
			ignored = append(ignored, strings.Fields(line[i+1:])...)
			continue
		}
	}
	return ignored
}

//TODO:take ignore groups into consideration?
