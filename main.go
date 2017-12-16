package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"strings"
)

func main() {
	args := os.Args[1:]

	if len(args) == 2 {
		address := args[1]
		if args[0] == "add" {
			_, err := addSite(address)
			if err != nil {
				fmt.Printf("Could not add %s. %s\n", address, err.Error())
			} else {
				fmt.Printf("%s added\n", address)
			}
		} else if args[0] == "remove" {
			_, err := removeSite(address)
			if err != nil {
				fmt.Printf("Could not remove %s. %s\n", address, err.Error())
			} else {
				fmt.Printf("%s removed\n", address)
			}
		}
	} else {
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

func getRcFileLocation() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s/%s", usr.HomeDir, ".stillup")
}

func addSite(site string) (bool, error) {

	urlObj, err := url.ParseRequestURI(site)

	if err != nil {
		return false, err
	}

	f, err := os.OpenFile(getRcFileLocation(), os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	_, writeError := f.WriteString(fmt.Sprintf("%s\n", urlObj.String()))
	if writeError != nil {
		return false, writeError
	}

	return true, nil
}

func removeSite(site string) (bool, error) {
	input, err := ioutil.ReadFile(getRcFileLocation())
	if err != nil {
		return false, err
	}

	lines := strings.Split(string(input), "\n")
	var newLines []string

	for _, line := range lines {
		if strings.TrimRight(line, "\n") != site {
			newLines = append(newLines, line)
		}
	}
	output := strings.Join(newLines, "\n")
	err = ioutil.WriteFile(getRcFileLocation(), []byte(output), 0644)
	if err != nil {
		return false, err
	}
	return true, nil
}
