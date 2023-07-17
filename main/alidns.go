package main

import (
	"fmt"
	"os"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/kdiot/alidns-console/console"
)

var (
	version = "0.0.1"
)

func Version() string {
	return version
}

func main() {

	usage := func() {
		fmt.Println("Usage:  alidns COMMAND")
		fmt.Println("Commands:")
		fmt.Println("  ls          List domain name records")
		fmt.Println("  add         Create a new domain name record")
		fmt.Println("  mod         Modify domain name record by RecordId")
		fmt.Println("  rm          Remove given domain name record by RecordId")
		fmt.Println("  ddns        Automatically update the domain name A record when a change in the external IP address is detected")
		fmt.Println("  help        Print help information for specific commands, such as: ls, add, etc.")
		fmt.Println("  version     Show the alidns version information")
	}

	if len(os.Args) < 2 {
		fmt.Println("Error! Not enough command line arguments.")
		usage()
		return
	}

	commands := map[string]console.Command{}
	if cmd := console.NewCmdLs(); cmd != nil {
		commands[cmd.Name()] = cmd
	}
	if cmd := console.NewCmdAdd(); cmd != nil {
		commands[cmd.Name()] = cmd
	}
	if cmd := console.NewCmdMod(); cmd != nil {
		commands[cmd.Name()] = cmd
	}
	if cmd := console.NewCmdRm(); cmd != nil {
		commands[cmd.Name()] = cmd
	}
	if cmd := console.NewCmdDdns(); cmd != nil {
		commands[cmd.Name()] = cmd
	}

	name := os.Args[1]
	if cmd, ok := commands[name]; ok {
		if err := cmd.Parse(os.Args[2:]); err != nil {
			fmt.Printf("Illegal parameter of command '%s': %s\n", name, err.Error())
			return
		}
		if err := cmd.Check(); err != nil {
			fmt.Printf("Command line '%s' parameter check failed: %s\n", name, err.Error())
			return
		}
		if err := cmd.Execute(); err != nil {
			var msg string
			if e, ok := err.(*tea.SDKError); ok {
				msg = fmt.Sprintf("ErrCode: %s, %s", *e.Code, *e.Message)
			} else {
				msg = err.Error()
			}
			fmt.Printf("Failed to execute '%s' command! [%s.]\n", name, msg)
		}
	} else if name == "help" {
		if len(os.Args) > 2 {
			cmd, ok := commands[os.Args[2]]
			if !ok {
				fmt.Printf("Error! Could not find help for command '%s'.\n", os.Args[2])
				usage()
				return
			}
			cmd.Usage()
		} else {
			usage()
		}
	} else if name == "version" {
		fmt.Println(Version())
	} else {
		usage()
	}
}
