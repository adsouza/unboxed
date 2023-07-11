package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func reportError(statusCode int, msg string, w http.ResponseWriter) {
	log.Print(msg)
	w.WriteHeader(statusCode)
	fmt.Fprintln(w, msg)
}

func reqHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	default:
		reportError(http.StatusMethodNotAllowed, fmt.Sprintf("Unsupported HTTP method %s.", r.Method), w)
	}
}

func main() {
	arg := os.Args[1]
	sql, err := os.ReadFile(arg + ".sql")
	cmd := exec.Command("sqlite3", "-header", "-html", "-readonly", "music.sqlite3")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	io.WriteString(stdin, string(sql))
	stdin.Close()
	if err := cmd.Run(); err != nil {
		log.Fatalf("SQL execution failed with error: %s\n", err)
	}
	output := string(stdout.Bytes())
	fmt.Printf("<!DOCTYPE html><html>\n<table>\n%s</table>\n</html>\n", output)
	http.HandleFunc("/", reqHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
