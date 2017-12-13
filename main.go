package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
)

func main() {
	file := getRc()
	defer file.Close()
	scanner := bufio.NewScanner(file)
	str := make(chan string)
	numSites := 0
	for scanner.Scan() {
		go checksite(scanner.Text(), str)
		numSites += 1
	}

	for i := 0; i < numSites; i++ {
		fmt.Print(<-str)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func getRc() *os.File {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	rcFile := fmt.Sprintf("%s/%s", usr.HomeDir, ".stillup")
	file, err := os.Open(rcFile)

	if err != nil {
		fmt.Printf("%s is not set.\nCreate a file .stillup in your homedir containing a list of sites", rcFile)
	}
	return file
}

func checksite(site string, str chan string) {
	urlstring, err := url.ParseRequestURI(site)
	if err != nil {
		str <- fmt.Sprintf("%s is not a valid URL\n\n", site)
	}

	resp, err := http.Get(site)

	if err != nil {
		str <- fmt.Sprintf("[404] %s\n", urlstring.Hostname())
	} else {
		str <- fmt.Sprintf("[%d] %s\n", resp.StatusCode, urlstring.Hostname())
	}
}
