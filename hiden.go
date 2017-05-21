package main

import (
	"fmt"
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
				fmt.Println("try install here")
				return nil
			},
		},
	}

	app.Run(os.Args)
}
