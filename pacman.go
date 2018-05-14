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
	"sync"
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
	var mux = &sync.Mutex{}
	var c = make(chan int)
	var toIgnore = make(chan string, 10)
	var visited = make(map[string]int)
	visited[confFile]++
	go readConf(confFile, visited, mux, toIgnore, c)
	for {
		select {
		case pkgGlob := <-toIgnore:
			*ignored = append(*ignored, pkgGlob)
		case <-c:
			ready <- 1
			return
		}
	}
}

func readConf(pacmanConf string, visited map[string]int, mux *sync.Mutex, toIgnore chan string, c chan int) {
	defer func(cc chan int) { cc <- 1 }(c)

	f, err := os.Open(pacmanConf)
	if err != nil {
		return
	}
	defer f.Close()
	conf := bufio.NewReader(f)

	var ready = make(chan int)
	var threads int

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
			file := strings.Fields(line[i+1:])[0]
			mux.Lock()
			if visited[file] == 0 {
				visited[file]++
				mux.Unlock()
				go readConf(file, visited, mux, toIgnore, ready)
				threads++
			} else {
				mux.Unlock()
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
				toIgnore <- pkgGlob
			}
			continue
		}
	}
	for threads > 0 {
		<-ready
		threads--
	}
}
