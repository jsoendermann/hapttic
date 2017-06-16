package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"unicode/utf8"
)

const version = "1.0.0"

var scriptFileName *string

type marshallableRequest struct {
	Method string
	URL    string
	Proto  string
	Host   string

	Header http.Header

	ContentLength int64
	Body          string
	Form          url.Values
	PostForm      url.Values
}

func init() {
	log.SetOutput(os.Stdout)
}

func handleFunc(w http.ResponseWriter, req *http.Request) {
	if _, err := os.Stat(*scriptFileName); os.IsNotExist(err) {
		log.Fatal("The request handling script " + *scriptFileName + " does not exist.")
	}

	bodyBuffer := new(bytes.Buffer)
	bodyBuffer.ReadFrom(req.Body)
	body := bodyBuffer.String()

	req.ParseForm()

	marshallableReq := marshallableRequest{
		Method: req.Method,
		URL:    req.URL.String(),
		Proto:  req.Proto,
		Host:   req.Host,

		Header: req.Header,

		ContentLength: req.ContentLength,
		Body:          body,
		Form:          req.Form,
		PostForm:      req.PostForm,
	}

	requestJson, err := json.Marshal(marshallableReq)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Executing " + *scriptFileName)

	out, err := exec.Command("/bin/bash", *scriptFileName, string(requestJson)).Output()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "500 Internal Server Error")
	} else {
		w.Write(out)
	}
}

func main() {
	printVersion := flag.Bool("v", false, "Print version and exit.")
	printUsage := flag.Bool("u", false, "Print usage and exit")
	host := flag.String("h", "", "The host to bind to, e.g. 0.0.0.0 or localhost.")
	port := flag.String("p", "8080", "The port to listen on.")
	userScriptFileName := flag.String("f", "./hapttic_request_handler.sh", "The script that is called to handle requests.")
	flag.Parse()

	if *printVersion {
		fmt.Fprintf(os.Stderr, version+"\n")
		os.Exit(0)
	}

	if *printUsage {
		fmt.Fprintf(os.Stderr, "Usage of hapttic:\n")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if utf8.RuneCountInString(*userScriptFileName) == 0 {
		log.Fatal("The path to the request handling script can not be empty.")
	}

	scriptFileName, err := filepath.Abs(*userScriptFileName)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handleFunc)

	addr := *host + ":" + *port
	log.Println("Thanks for using hapttic")
	log.Println(fmt.Sprintf("Listening on %s", addr))
	log.Println(fmt.Sprintf("Forwarding requests to %s", scriptFileName))
	log.Fatal(http.ListenAndServe(addr, nil))
}
