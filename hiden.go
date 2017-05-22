package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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

func github_binary_install(name string) error {
	name = "https://api.github.com/repos/prometheus/prometheus/releases/latest"
	var release GithubRelease

	resp, err := http.Get(name)
	if err != nil {
		fmt.Printf("Error getting %s:\n\t%s\n", name, err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error getting text body in %s:\n\t%s\n", name, err)
		return err
	}

	err = json.Unmarshal(body, &release)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON in %s:\n\t%s\n ", name, err)
		return err
	}

	return release
}
