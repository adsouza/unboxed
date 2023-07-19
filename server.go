package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func reportError(statusCode int, msg string, w http.ResponseWriter) {
	log.Print(msg)
	w.WriteHeader(statusCode)
	fmt.Fprintln(w, msg)
}

func reqHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		page := strings.TrimPrefix(r.URL.Path, "/")
		switch page {
		case "":
			page = "all"
		case "top5-musicians":
			page = "top5-musicians-by-total-sales"
		case "greatest-hits":
			page = "songs-with-sales-over-20m"
		case "odd-years":
			page = "hits-from-odd-numbered-years"
		case "2nd-person":
			page = "songs-whose-title-contains-you"
		case "favicon.ico":
			w.WriteHeader(http.StatusNotFound)
			return
		default:
			reportError(http.StatusNotFound, fmt.Sprintf("Unrecognized path: %s.", page), w)
			return
		}
		fmt.Fprintln(w, render(page))
	default:
		reportError(http.StatusMethodNotAllowed, fmt.Sprintf("Unsupported HTTP method %s.", r.Method), w)
	}
}

func render(arg string) string {
	sql, err := os.ReadFile(arg + ".sql")
	if err != nil {
		log.Fatal(err)
	}
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
	results := string(stdout.Bytes())
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<table style="float:left">
%s</table>

<div style="float:left; margin-left:2em;">
<h1>Show only:</h1>
<ul style="line-height:1.5">
	<li><a href="/top5-musicians">The top five musicians by total sales</a></li>
	<li><a href="/greatest-hits">Songs with sales over $20 million</a></li>
	<li><a href="/odd-years">Songs released during odd-numbered years</a></li>
	<li><a href="/2nd-person">Songs whose title contains the word "you"</a></li>
</ul>
</div>
</html>
`,
		results)
}

func main() {
	fmt.Println(render("all"))
	http.HandleFunc("/", reqHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
