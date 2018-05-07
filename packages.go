package main

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
)

type pkgInfo struct {
	version    string
	aurVersion string
}

// GET PACKAGES
func getForeignPackages() (map[string]pkgInfo, error) {
	cmd := exec.Command("/bin/pacman", "-Qm")
	out, err := cmd.Output()

	if err != nil {
		return nil, errors.New("can't read pacman's output")
	}

	packages := make(map[string]pkgInfo)
	scanner := bufio.NewScanner(bytes.NewReader(out))
	scanner.Split(bufio.ScanWords)

	for i, p := 0, ""; scanner.Scan(); i++ {
		if i%2 == 0 {
			p = scanner.Text()
		} else {
			packages[p] = pkgInfo{version: scanner.Text()}
		}
	}
	return packages, nil
}
