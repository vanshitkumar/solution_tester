/*
Copyright © 2025 Vanshit hello@vanshit.me
*/
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func formatHeading(heading string) string {
	border := strings.Repeat("#", len(heading)+8)
	return fmt.Sprintf("%s\n### %s ###\n%s\n", border, heading, border)
}

func formatCompilerOutput(output []byte) string {
	if len(output) > 0 {
		combinedOutput := formatHeading("Compiler Output")
		combinedOutput += string(output)
		return combinedOutput
	}
	return ""
}

func compileFile(fileName string) error {
	command := fmt.Sprintf(COMPILE_COMMAND, fileName, "temp")
	parts := strings.Fields(command)
	c := exec.Command(parts[0], parts[1:]...)

	output, err := c.CombinedOutput()
	fmt.Print(formatCompilerOutput(output))
	return err
}

type Test struct{
	Input string `json:"input"`
	Output string `json:"output"`	
}

type Question struct {
	Name 	  string   `json:"name"`
	Group 	  string   `json:"group"`
	Url 	 string   `json:"url"`
	Interactive bool     `json:"interactive"` // optional
	TimeLimit   int      `json:"timeLimit"`
	MemoryLimit int      `json:"memoryLimit"`
	Tests       []Test   `json:"tests"`	
}


/// OUTPUT FORMATTING GUIDELINES
/// Each output ends with a newline
/// Only Heading are printed with ### only important messages
/// Test case results are printed with ✅ or ❌
/// In case of runtime error the output of the program is printed first followed by the error message
/// In case of compilation error the output of the compiler is printed

// fetches the question
// print if no question or interactive question
func getQuestion() *Question {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%s/", PORT))
	if err != nil || resp.StatusCode != 200 {
		fmt.Print(formatHeading("No Question"))
		return nil
	}

	question := Question{}
	err = json.NewDecoder(resp.Body).Decode(&question)
	if err != nil {
		fmt.Print(formatHeading("Failed to parse question"))
		return nil
	}

	if question.Interactive{
		fmt.Print(formatHeading("Interactive Question"))
	}
	return &question
}

// returns output and error if any
// timlimint -1 means no limit else in milliseconds
// prints nothing
func runTestCase(input string, timeLimit int) (string, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	if(timeLimit == -1){
		ctx = context.Background()
	}else{
		ctx, cancel = context.WithTimeout(context.Background(),  time.Millisecond * time.Duration(timeLimit))
		defer cancel()
	}

	c := exec.CommandContext(ctx, "./temp")
	c.Stdin = strings.NewReader(input)
	output, err := c.CombinedOutput()
	return string(output), err
}

func isWhitespace(c byte) bool {
    return c == ' ' || c == '\n' || c == '\t' || c == '\r'
}

// compares character by character ignoring whitespaces in actual output
// prints if test passes or fails and returns bool
func compareOutput(expected, actual string, tcCount int) bool {	
	i, j := 0, 0
	for i < len(expected) && j < len(actual) {
		if expected[i] == actual[j] {
			i++
			j++
		} else if isWhitespace(actual[j]) {
			j++
		} else {
			fmt.Printf("❌ TEST CASE %d FAILED\n", tcCount)
			return false
		}
	}

	
	for j < len(actual) {
		if isWhitespace(actual[j]) {
			j++
		} else {
			fmt.Printf("❌ TEST CASE %d FAILED\n", tcCount)
			return false
		}
	}

	for i < len(expected) {
		if isWhitespace(expected[i]) {
			i++
		} else {
			fmt.Printf("❌ TEST CASE %d FAILED\n", tcCount)
			return false
		}
	}

	fmt.Printf("✅ TEST CASE %d PASSED\n", tcCount)
	return true
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "submit your file for testing",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := compileFile(args[0])
		if err != nil {
			fmt.Print(formatHeading("Compilation Failed"))
			return
		}

		// here we guranteedly have a compiled file named temp


		question := getQuestion()
		if(question != nil){
			fmt.Printf("# %s\n# %s\n# %s\n", question.Group, question.Name, question.Url)
		}

		if question == nil || question.Interactive {
			// in this case we should just run the file and forward stdin and stdout
			// all the other cases we send test cases ourselves
			c := exec.Command("./temp")
			c.Stdin = cmd.InOrStdin()
			c.Stdout = cmd.OutOrStdout()
			c.Stderr = cmd.OutOrStderr()
			err := c.Run()
			if err != nil {
				fmt.Println("!!!!!!!! Runtime Error")
				fmt.Println(err.Error())
			}
			return
		}else{
			passed := 0
			for i, test := range question.Tests {
				output, err := runTestCase(test.Input, question.TimeLimit)
				if err != nil {
					fmt.Print(output)
					fmt.Println("!!!!!!!! Runtime Error")
					fmt.Println(err.Error())
					continue
				}
				if compareOutput(test.Output, output, i+1) {
					passed++
				}else{
					fmt.Println("Input:")
					fmt.Print(test.Input)
					fmt.Println("Expected Output:")
					fmt.Print(test.Output)
					fmt.Println("Actual Output:")
					fmt.Print(output)
				}
			}

			fmt.Printf("# Test Cases Passed: %d/%d\n", passed, len(question.Tests))
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
