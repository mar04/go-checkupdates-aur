package main

import (
	"bufio"
	"bytes"
	"fmt"
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
		panic(err)
	}
	result, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		panic(err)
	}
	return result < 0
}

func getIgnored(confFile string, ignored *[]string, ready chan int) {
	var visited = make(map[string]int)
	visited[confFile]++
	readConf(confFile, visited, ignored)
	ready <- 1
}

func readConf(pacmanConf string, visited map[string]int, ignored *[]string) {
	f, err := os.Open(pacmanConf)
	if err != nil {
		return
	}
	defer f.Close()
	scan := bufio.NewScanner(f)

	for scan.Scan() {
		line := scan.Text()
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "Include") {
			//strip comments
			if i := strings.IndexByte(line, '#'); i >= 0 {
				line = line[:i]
			}

			i := strings.IndexByte(line, '=')
			file := strings.Fields(line[i+1:])[0]
			if visited[file] == 0 {
				visited[file]++
				readConf(file, visited, ignored)
			}
			continue
		}
		if strings.HasPrefix(line, "IgnorePkg") {
			//strip comments
			if i := strings.IndexByte(line, '#'); i >= 0 {
				line = line[:i]
			}

			i := strings.IndexByte(line, '=')
			for _, pkgGlob := range strings.Fields(line[i+1:]) {
				*ignored = append(*ignored, pkgGlob)
			}
			continue
		}
	}
}
