package main

import (
	"encoding/json"
	"fmt"
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
				fmt.Printf("First argument: %s\n", name)
				github_binary_install(name)
				return nil
			},
		},
	}

	app.Run(os.Args)
}

func github_binary_install(name string) error {
	name = "https://api.github.com/repos/prometheus/prometheus/releases/latest"

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

	err := json.Unmarshal()
	if err != nil {
		fmt.Printf("Error unmarshalling JSON in %s:\n\t%s\n ", name, err)
		return err
	}
	return nil
}
