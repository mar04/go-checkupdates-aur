package main

import (
	"testing"
)

func Test_constructQuery(t *testing.T) {
	type args struct {
		packages map[string]*pkgInfo
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Map 'packages' is nil",
			args: args{packages: nil},
			want: "https://aur.archlinux.org/rpc/?v=5&type=info",
		},
		{
			name: "Map 'packages' is empty",
			args: args{packages: make(map[string]*pkgInfo)},
			want: "https://aur.archlinux.org/rpc/?v=5&type=info",
		},
		{
			name: "Map has single entry",
			args: args{packages: map[string]*pkgInfo{
				"yaourt": nil,
			}},
			want: "https://aur.archlinux.org/rpc/?v=5&type=info&arg[]=yaourt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := constructQuery(tt.args.packages); got != tt.want {
				t.Errorf("constructQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getAurVersions(t *testing.T) {
	type args struct {
		packages map[string]*pkgInfo
	}
	tests := []struct {
		name string
		args args
	}{
		//no network
		//error response
		//malformed json
		//pkginfo is nil
		{
			name: "Map 'packages' is nil",
			args: args{ packages: nil},
		},
		{
			name: "Map 'packages' is empty",
			args: args{packages: make(map[string]*pkgInfo)},
		},
		{
			name: "Map has single entry",
			args: args{packages: map[string]*pkgInfo{
				"yaourt": nil,
			}},
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getAurVersions(tt.args.packages)
		})
	}
}
