package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Makepad-fr/semver/cli/cmd"
	"github.com/Makepad-fr/semver/semver"
)

func init() {
	// Add diff command
	cmd.AddCommand(cmd.DiffCommand)
}

func Run() error {
	version := os.Args[1]
	sv, err := semver.Parse(version)
	if err != nil {
		return err
	}

	verbosePtr := flag.Bool("verbose", false, "Show verbose output")
	flag.Parse()
	if verbosePtr != nil && *verbosePtr {
		fmt.Printf("Major %s\n", sv.Major)
		fmt.Printf("Minor %s\n", sv.Minor)
		fmt.Printf("Patch %s\n", sv.Patch)
		if sv.PreRelease != nil {
			fmt.Printf("PreRelease %s\n", *sv.PreRelease)
		}
		if sv.BuildMetaData != nil {
			fmt.Printf("BuildMetaData %s", *sv.BuildMetaData)
		}
	}
	return nil
}

func main() {
	filename := os.Args[0]
	if len(os.Args) > 1 {
		subCommand := os.Args[1]

		cmd, ok := cmd.GetCommand(subCommand)
		if ok {
			// Sub command exists
			os.Args = os.Args[1:]
			err := cmd.Run()
			if err != nil {
				fmt.Printf("Error while running command %s %s : %v\n", filename, subCommand, err)
				flag.Usage()
				os.Exit(1)
			}
			return
		}
		err := Run()
		if err != nil {
			fmt.Println(err)
			flag.Usage()
			os.Exit(1)
		}
		return
	}
	flag.Usage()
}
