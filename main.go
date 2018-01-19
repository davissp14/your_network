package main

import (
	"fmt"
	"mydb/configuration"
	"mydb/server"
	"os"

	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	cli.AppHelpTemplate = cliTemplate()
	app.Name = "YourNetwork"
	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "generate configuration file `config.json` (Required)",
			Action: func(c *cli.Context) error {
				configuration.Init()
				os.Exit(0)
				return nil
			},
		},
		{
			Name:  "start",
			Usage: "Starts your server",
			Subcommands: []cli.Command{
				{
					Name:  "config",
					Usage: "Specify the path to your `config.json` file.",
					Action: func(c *cli.Context) error {
						config, err := configuration.Load(c.Args().First())
						if err != nil {
							fmt.Println(err)
							os.Exit(1)
						}
						server.Start(config)
						return nil
					},
				},
			},
		},
	}

	app.Run(os.Args)
}

func cliTemplate() string {
	return fmt.Sprintf(`{{if .VisibleCommands}}COMMANDS:{{range .VisibleCategories}}{{if .Name}}
	 {{.Name}}:{{end}}{{range .VisibleCommands}}
		 {{join .Names ", "}}{{"\t"}}{{end}}{{end}}{{end}}
	`)
}
