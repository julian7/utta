package main

import (
	"fmt"
	"os"
	"os/exec"
)

const VERFILE = "version.txt"

func main() {
	verfile, err := os.OpenFile(VERFILE, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Printf("error: cannot open %s: %v\n", VERFILE, err)
		os.Exit(1)
	}
	defer verfile.Close()

	cmd := exec.Command("git", "describe", "--tags", "--always", "--dirty")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("error: error getting output of git describe: %v\n", err)
		os.Exit(1)
	}

	_, err = verfile.Write(output)
	if err != nil {
		fmt.Printf("error: error getting output of git describe: %v\n", err)
		os.Exit(1)
	}
}
