package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	//"runtime"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "hiden"
	app.Usage = "install and manage binary files from the Internet"

	app.Commands = []cli.Command{
		{
			Name:    "install",
			Aliases: []string{"i"},
			Usage:   "install new package",
			Action: func(c *cli.Context) error {
				name := c.Args().Get(0)
				github_binary_install(name)
				return nil
			},
		},
	}

	app.Run(os.Args)
}

type GithubRelease struct {
	Id        int           `json:"id"`
	TagName   string        `json:"tag_name"`
	Url       string        `json:"url"`
	CreatedAt string        `json:"created_at"`
	Assets    []GithubAsset `json:"assets"`
}

type GithubAsset struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Url       string `json:"browser_download_url"`
}

func (release GithubRelease) BestAsset() (GithubAsset, error) {
	var asset GithubAsset
	for _, elem := range release.Assets {
		os, arch := elem.Type()
		if os == "darwin" && arch == "amd64" {
			asset = elem
		}
	}
	return asset, nil
}

var os_mac = regexp.MustCompile(`darwin|mac|osx`)
var os_lin = regexp.MustCompile(`linux`)
var os_win = regexp.MustCompile(`win|ms`)
var os_bsd = regexp.MustCompile(`bsd`)
var os_sol = regexp.MustCompile(`solaris`)

var arch_amd64 = regexp.MustCompile(`(amd|x(86[_-])?)64|64[^a-zA-Z0-9]`)
var arch_arm = regexp.MustCompile(`arm`)
var arch_powerpc = regexp.MustCompile(`powerpc|ppc`)

func match(text string, re *regexp.Regexp) bool {
	if re.FindStringIndex(text) == nil {
		return false
	}
	return true
}

func (asset GithubAsset) Type() (string, string) {
	os := "unknown"
	arch := "i386"

	// FIXME: dumb regex check
	//        this may check wrong in so many cases
	switch {
	case match(asset.Name, os_mac):
		os = "darwin"
	case match(asset.Name, os_lin):
		os = "linux"
	case match(asset.Name, os_win):
		os = "windows"
	case match(asset.Name, os_bsd):
		os = "bsd"
	case match(asset.Name, os_sol):
		os = "solaris"
	}

	switch {
	case match(asset.Name, arch_amd64):
		arch = "amd64"
	case match(asset.Name, arch_powerpc):
		arch = "powerpc"
	case match(asset.Name, arch_arm):
		arch = "arm"
	}

	return os, arch
}

func github_binary_install(name string) error {
	releases, _ := get_latest_release(name)
	asset, _ := releases.BestAsset()
	fmt.Printf("GOOOGOOOGOOGOO [%s]\n", asset.Url)
	return nil
}

func get_latest_release(name string) (GithubRelease, error) {
	name = "https://api.github.com/repos/prometheus/prometheus/releases/latest"
	var release GithubRelease

	resp, err := http.Get(name)
	if err != nil {
		fmt.Printf("Error getting %s:\n\t%s\n", name, err)
		return release, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error getting text body in %s:\n\t%s\n", name, err)
		return release, err
	}

	err = json.Unmarshal(body, &release)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON in %s:\n\t%s\n ", name, err)
		return release, err
	}

	return release, nil
}
