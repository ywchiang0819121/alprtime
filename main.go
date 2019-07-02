package main

import (
	"bufio"
    "fmt"
    "log"
	"os"
	"strings"
)

func main() {
    file, err := os.Open("/var/log/iptables.log")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		list := strings.Split(line, " ")
		for i, substr := range list {
			println(substr)
		}
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
	}
}