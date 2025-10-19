/*
Copyright © 2025 Vanshit hello@vanshit.me
*/
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func getInTxt() (bool, string) {
	data, err := os.ReadFile("in.txt")
	if err != nil {
		return false, ""
	}
	return true, string(data)
}

const (
	IN_METHOD_STDIN = iota
	IN_METHOD_FILE
	IN_METHOD_SAMPLES
)

var (
	inTxtContents string
)

func printInputMethod(q *Question) int {
	var hasInTxt bool
	hasInTxt, inTxtContents = getInTxt()
	if hasInTxt {
		printNormal("Input Method: 'in.txt'")
		return IN_METHOD_FILE
	} else if q != nil && q.Interactive {
		printNormal("Input Method: Stdin (Interactive Question)")
		return IN_METHOD_STDIN
	} else if q != nil {
		printNormal("Input Method: Sample Test Cases")
		return IN_METHOD_SAMPLES
	} else {
		printNormal("Input Method: Stdin")
		return IN_METHOD_STDIN
	}
}

func compileFile(fileName string) error {
	command := fmt.Sprintf(COMPILE_COMMAND, fileName, "temp")
	parts := strings.Fields(command)
	c := exec.Command(parts[0], parts[1:]...)

	output, err := c.CombinedOutput()
	printError("Compiler Output", string(output))
	return err
}

type Test struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type Question struct {
	Name        string `json:"name"`
	Group       string `json:"group"`
	Url         string `json:"url"`
	Interactive bool   `json:"interactive"` // optional
	TimeLimit   int    `json:"timeLimit"`
	MemoryLimit int    `json:"memoryLimit"`
	Tests       []Test `json:"tests"`
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
		return nil
	}

	question := Question{}
	err = json.NewDecoder(resp.Body).Decode(&question)
	if err != nil {
		return nil
	}
	return &question
}

// returns output and error if any
// timlimint -1 means no limit else in milliseconds
// prints nothing
func runTestCase(input string, timeLimit int) (string, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	if timeLimit == -1 {
		ctx = context.Background()
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*time.Duration(timeLimit))
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
func compareOutput(expected, actual string) bool {
	i, j := 0, 0
	for i < len(expected) && j < len(actual) {
		if expected[i] == actual[j] {
			i++
			j++
		} else if isWhitespace(actual[j]) {
			j++
		} else {
			return false
		}
	}

	for j < len(actual) {
		if isWhitespace(actual[j]) {
			j++
		} else {
			return false
		}
	}

	for i < len(expected) {
		if isWhitespace(expected[i]) {
			i++
		} else {
			return false
		}
	}

	return true
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "submit your file for testing",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		question := getQuestion()
		printQuestion(question)

		err := compileFile(args[0])
		if err != nil {
			printRed("Compilation Failed")
			return
		} else {
			printGreen("Compilation Successfull")
		}

		method := printInputMethod(question)

		switch method {
		case IN_METHOD_STDIN:
			c := exec.Command("./temp")
			c.Stdin = cmd.InOrStdin()
			c.Stdout = cmd.OutOrStdout()
			c.Stderr = cmd.OutOrStderr()
			err := c.Run()
			if err != nil {
				printError("Runtime Error", err.Error());
			}
		case IN_METHOD_FILE:
			c := exec.Command("./temp")
			c.Stdin = strings.NewReader(inTxtContents)
			c.Stdout = cmd.OutOrStdout()
			c.Stderr = cmd.OutOrStderr()
			err := c.Run()
			if err != nil {
				printError("Runtime Error", err.Error())
			}
		case IN_METHOD_SAMPLES:
			passed := 0
			for i, test := range question.Tests {
				output, err := runTestCase(test.Input, question.TimeLimit)
				if err != nil {
					printRed(fmt.Sprintf("Test Case %d failed", i+1))
					printError("Runtime Error", err.Error())
					fmt.Println("#  Output:")
					fmt.Print(output)
					continue
				}
				if compareOutput(test.Output, output) {
					printGreen(fmt.Sprintf("Test Case %d passed", i+1))
					passed++
				} else {
					printRed(fmt.Sprintf("Test Case %d failed", i+1))
					fmt.Println("\n#  Input:")
					fmt.Print(test.Input)
					fmt.Println("\n#  Expected Output:")
					fmt.Print(test.Output)
					fmt.Println("\n#  Actual Output:")
					fmt.Print(output)
				}
			}

			if(passed == len(question.Tests)){
				printGreen("All Test Cases Passed ✅")
			} else {
				printRed(fmt.Sprintf("Passed %d out of %d test cases", passed, len(question.Tests)))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
