package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strconv"
    "strings"
    "time"
    "os/exec"
    // "net/http"
)

const ntpPeriodMAX = 5

func getLastBootTime() string {
    out, err := exec.Command("who", "-b").Output()
    if err != nil {
            panic(err)
    }
    t := strings.TrimSpace(string(out))
    t = strings.TrimPrefix(t, "system boot")
    t = strings.TrimSpace(t)
    return t
}

func getTimezone() string {
    out, err := exec.Command("date", "+%Z").Output()
    if err != nil {
            panic(err)
    }
    return strings.TrimSpace(string(out))
}

func getLastSystemBootTime() (time.Time, error) {
    return time.Parse(`2006-01-02 15:04MST`, getLastBootTime()+getTimezone())
}

func main() {
    file, err := os.Open("iptables.log")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    var watchList []string
	watchList = append(watchList, "140.118.127.164")
    watchList = append(watchList, "192.168.80.1")	
    scanner := bufio.NewScanner(file)
    ipWithLastVisit := make(map[string]time.Duration)

    for scanner.Scan() {
        line := scanner.Text()
        list := strings.Split(line, " ")
        epch := strings.Replace(list[6][1:], "]", "", -1)
        epchfloat, _ := strconv.ParseFloat(epch, 64)
        key := string(list[10][4:])

        lastEntry := time.Unix(int64(epchfloat)+time.Now().Unix(), int64(100000*(epchfloat-float64(int64(epchfloat)))))
        now := time.Now()

        diff := now.Sub(lastEntry)

        ipWithLastVisit[key] = diff
    }

    for ip := range watchList {
        if ipWithLastVisit[watchList[ip]].Minutes() > ntpPeriodMAX {
            fmt.Print("ip : ", watchList[ip], "\tlast entry was ", int(ipWithLastVisit[watchList[ip]].Hours()), "\tHours\t", int(ipWithLastVisit[watchList[ip]].Minutes())%60, "\tMinutes\t", int(ipWithLastVisit[watchList[ip]].Seconds())%60, "\tSeconds ago")
            fmt.Println(" alert!")
        } else {
            fmt.Println()
        }
    }
    
    

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}
