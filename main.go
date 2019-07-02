package main

import (
	"bufio"
    "fmt"
    "log"
    "os"
)

func main() {
    file, err := os.Open("/var/log/iptables.log")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

	scanner := bufio.NewScanner(file)
	cnt := 1
    for scanner.Scan() {
		fmt.Println(scanner.Text())
		cnt++
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
	}
}