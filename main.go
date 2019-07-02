package main

import (
	"bufio"
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
	var ipWithLastVisit map[string]string


	for scanner.Scan() {
		line := scanner.Text()
		list := strings.Split(line, " ")
		// for _, substr := range list {
		// 	println(substr)
		// }
		key := string(list[10][4:])
		val := string(list[0] + "-" + list[2] + "-"  + list[3])

		ipWithLastVisit[key] = val
        println(key)
        println(ipWithLastVisit[key])
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
	}
}