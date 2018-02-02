package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"your_network/client"
	"your_network/configuration"
)

func main() {
	// Subcommands
	serverCommand := flag.NewFlagSet("server", flag.ExitOnError)
	// initCommand := flag.NewFlagSet("init", flag.ExitOnError)
	addNodeCommand := flag.NewFlagSet("add_node", flag.ExitOnError)
	listNodesCommand := flag.NewFlagSet("list_nodes", flag.ExitOnError)
	// pingCommand := flag.NewFlagSet("ping", flag.ExitOnError)
	// listFilesCommand := flag.NewFlagSet("list_files", flag.ExitOnError)
	// downloadCommand := flag.NewFlagSet("download", flag.ExitOnError)

	// Server subcommands
	serverConfigPtr := serverCommand.String("config", "./config.json", "Configuration file (Required)")
	// Add node subcommands
	addNodeHostnamePtr := addNodeCommand.String("node", "", "Target node. e.g. '192.168.1.4:4000'. (Required)")
	// Ping subcommands
	// pingHostnamePtr := pingCommand.String("node", "", "Target node. e.g. '192.168.1.4:4000'. (Required)")
	// // List files subcommmands
	// listFilesHostnamePtr := listFilesCommand.String("node", "", "Target node. '192.168.1.4:4000'. (Required)")
	// // Download subcommands
	// downloadHostnamePtr := downloadCommand.String("node", "", "Target node. '192.168.1.4:4000'. (Required)")
	// downloadFilenamePtr := downloadCommand.String("filename", "", "Filename you want to download. (Required)")

	if len(os.Args) < 2 {
		fmt.Printf(
			`Your Network

Available Commands
==============
server:     Starts server.
init:       Generates config file.
add_node:   Add Remote / Local node to your network.
list_nodes: List all nodes attached to your network.
list_files: List all shared files on a given node.
download:   Download remote file.
ping:       Ping Remote / Local node within your network.
`)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "server":
		serverCommand.Parse(os.Args[2:])
	case "init":
		err := configuration.Init()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	case "add_node":
		addNodeCommand.Parse(os.Args[2:])
	case "list_nodes":
		listNodesCommand.Parse(os.Args[2:])
	// case "list_files":
	// 	listFilesCommand.Parse(os.Args[2:])
	// case "download":
	// 	downloadCommand.Parse(os.Args[2:])
	// case "ping":
	// 	pingCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	config, err := configuration.Load(*serverConfigPtr)
	if err != nil {
		fmt.Printf("Unable to find config.json file at `%s`\n", *serverConfigPtr)
		os.Exit(1)
	}

	cl := client.Client{Config: config}
	// Check which subcommand was Parsed using the FlagSet.Parsed() function. Handle each case accordingly.
	// FlagSet.Parse() will evaluate to false if no flags were parsed (i.e. the user did not provide any flags)
	if serverCommand.Parsed() {
		// Required Flags
		if *serverConfigPtr == "" {
			serverCommand.PrintDefaults()
			os.Exit(1)
		}
		cl.StartServer()
	}

	if addNodeCommand.Parsed() {
		if *addNodeHostnamePtr == "" {
			addNodeCommand.PrintDefaults()
			os.Exit(1)
		}
		cl.AddNode(strings.TrimSpace(*addNodeHostnamePtr))
	}

	if listNodesCommand.Parsed() {
		cl.ListNodes()
	}
	//
	// if listFilesCommand.Parsed() {
	// 	if *listFilesHostnamePtr == "" {
	// 		listFilesCommand.PrintDefaults()
	// 		os.Exit(1)
	// 	}
	// 	cl.ListFiles(strings.TrimSpace(*listFilesHostnamePtr))
	// }
	//
	// if downloadCommand.Parsed() {
	// 	if *downloadFilenamePtr == "" || *downloadHostnamePtr == "" {
	// 		downloadCommand.PrintDefaults()
	// 		os.Exit(1)
	// 	}
	// 	cl.Download(
	// 		strings.TrimSpace(*downloadHostnamePtr),
	// 		strings.TrimSpace(*downloadFilenamePtr),
	// 	)
	// }
	//
	// if pingCommand.Parsed() {
	// 	if *pingHostnamePtr == "" {
	// 		pingCommand.PrintDefaults()
	// 		os.Exit(1)
	// 	}
	// 	cl.Ping(*pingHostnamePtr)
	// }

}

// app := cli.NewApp()
// cli.AppHelpTemplate = cliTemplate()
// app.Name = "YourNetwork"
// app.Commands = []cli.Command{
// 	{
// 		Name:  "init",
// 		Usage: "generate configuration file `config.json` (Required)",
// 		Action: func(c *cli.Context) error {
// 			configuration.Init()
// 			os.Exit(0)
// 			return nil
// 		},
// 	},
// 	{
// 		Name:  "start",
// 		Usage: "Starts your server",
// 		Subcommands: []cli.Command{
// 			{
// 				Name:  "config",
// 				Usage: "Specify the path to your `config.json` file.",
// 				Action: func(c *cli.Context) error {
// 					config, err := configuration.Load(c.Args().First())
// 					if err != nil {
// 						fmt.Println(err)
// 						os.Exit(1)
// 					}
// 					server.Start(config)
// 					return nil
// 				},
// 			},
// 		},
// 	},
// }
//
// app.Run(os.Args)
// }

// func cliTemplate() string {
// 	return fmt.Sprintf(`{{if .VisibleCommands}}COMMANDS:{{range .VisibleCategories}}{{if .Name}}
// 	 {{.Name}}:{{end}}{{range .VisibleCommands}}
// 		 {{join .Names ", "}}{{"\t"}}{{end}}{{end}}{{end}}
// 	`)
// }
