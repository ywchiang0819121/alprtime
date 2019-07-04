package main

import (
	"bufio"
	"fmt"
    "log"
	"os"
	"strings"
	"time"
	"strconv"
)

func strtoMonth(str string) int{
	switch str {
	case "Jan":
		return 1
	case "Feb":
		return 2
	case "Mar":
		return 3
	case "Apr":
		return 4
	case "May":
		return 5
	case "Jun":
		return 6
	case "Jul":
		return 7
	case "Aug":
		return 8
	case "Sep":
		return 9
	case "Oct":
		return 10
	case "Nov":
		return 11
	case "Dec":
		return 12
	}
	fmt.Println("not found Month : " + str)
	return 1
}

const ntpPeriodMAX = 5

func main() {
    file, err := os.Open("iptables.log")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

	var watchList[] string
	watchList = append(watchList, "140.118.127.164")
	scanner := bufio.NewScanner(file)
	ipWithLastVisit := make(map[string]time.Duration)
	tname, toffset := time.Now().Zone()
	tyear := time.Now().Year()

	for scanner.Scan() {
		line := scanner.Text()
		list := strings.Split(line, " ")
		key := string(list[10][4:])
		// val := string(list[0] + "-" + list[2] + "-"  + list[3])
		hms := strings.Split(string(list[3]), ":")
		d, _ := strconv.Atoi(list[2])
		h, _ := strconv.Atoi(hms[0])
		m, _ := strconv.Atoi(hms[1])
		s, _ := strconv.Atoi(hms[2])
		
		lastEntry := time.Date(tyear, time.Month(strtoMonth(list[0])), d, h, m, s,0, time.FixedZone(tname, toffset))
		now := time.Now()

		diff := now.Sub(lastEntry)

		ipWithLastVisit[key] = diff
	}
	
	for ip := range watchList {
		if ipWithLastVisit[watchList[ip]].Minutes() > ntpPeriodMAX{
			fmt.Print("ip : ", watchList[ip], "\tlast entry was ", int(ipWithLastVisit[watchList[ip]].Hours()), "\tHours\t", int(ipWithLastVisit[watchList[ip]].Minutes()) % 60, "\tMinutes\t", int(ipWithLastVisit[watchList[ip]].Seconds()) % 60, "\tSeconds ago")
			fmt.Println(" alert!")
		} else {
			fmt.Println()
		}
	}



    if err := scanner.Err(); err != nil {
        log.Fatal(err)
	}
}