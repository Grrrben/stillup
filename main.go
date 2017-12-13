package main

import (
	"fmt"
	"net/http"
	"os/user"
	"log"
	"os"
	"bufio"
	"net/url"
)

func main() {
	file := getRc();
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		success := checksite(scanner.Text())
		if (!success) {
			fmt.Printf("%s is not a valid URL\n\n", scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func getRc() *os.File {
	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
	}

	rcFile := fmt.Sprintf("%s/%s", usr.HomeDir,  ".stillup")
	file, err := os.Open(rcFile)

	if err != nil {
		fmt.Printf("%s is not set.\nCreate a file .stillup in your homedir containing a list of sites", rcFile)
	}
	return file
}

func checksite (site string) bool {
	urlstring, err := url.ParseRequestURI(site)
	if err != nil {
		return false
	}

	resp, err := http.Get(site)

	if err != nil {
		fmt.Printf("[404] %s\n", urlstring.Hostname())
	} else {
		fmt.Printf("[%d] %s\n", resp.StatusCode, urlstring.Hostname())
	}
	return true
}
