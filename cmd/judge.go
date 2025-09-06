/*
Copyright © 2025 Vanshit hello@vanshit.me
*/
package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/spf13/cobra"
)

var judgeCmd = &cobra.Command{
	Use:   "judge",
	Short: "starts the judge in the local console",
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

var (
	question     []byte
	questionMutex sync.Mutex
)

func start() {

	http.HandleFunc("POST /{$}", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		questionMutex.Lock()
		defer questionMutex.Unlock()
		question = bytes.Clone(body)
	})

	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		questionMutex.Lock()
		defer questionMutex.Unlock()
		if question == nil {
			http.Error(w, "No question available", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(question)
	})
	
	http.ListenAndServe(fmt.Sprintf(":%s", PORT), http.DefaultServeMux)
}

func init() {
	rootCmd.AddCommand(judgeCmd)
}
