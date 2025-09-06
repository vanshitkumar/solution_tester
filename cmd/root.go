/*
Copyright © 2025 Vanshit hello@vanshit.me
*/
package main

import (
	"os"

	"github.com/spf13/cobra"
)



var rootCmd = &cobra.Command{
	Use:   "solution_tester",
	Short: "This is a cli tool for testing your solutions against testcases fetched from various online judges(Codeforces, CodeChef, etc.).",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}