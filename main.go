package main

import (
    "bufio"
    "log"
    "os"
    "strconv"
    "strings"
    "time"
	"os/exec"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"bytes"
)


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

func getLastSystemBootTime() (time.Time, error) {
	// "2019-08-11 18:59"
    // return time.ParseInLocation(`2006-01-02 15:04`, getLastBootTime(), time.Local)
    return time.ParseInLocation(`2006-01-02 15:04`, "2019-08-11 18:59", time.Local)
}

type Device struct {
	Ip string 					`json:"serial_number"`
	Ipcams_id int				`json:"ipcams_id"`
	Recognizing_server_id int	`json:"recognizing_servers_id"`
	Time string					`json:"time"`
}

type Log struct {
	Token		string		`json:"token"`
	Server_time	string		`json:"server_time"`
	Logs		[]Device	`json:"logs"`
}

func main() {
	var ntplog Log
	ntplog.Token = "Fz7Brl6TGuhe3dI4B3Wfk1cXp4oqua44LrEKW52juOYiU32NOkdnWfJxQ3pvr9F7IoX3wAL6xeMKWcTQcPtxgjt5lV8a23U96zqU"
	ntplog.Server_time = strings.Split(time.Now().String(), " +")[0]
	byteStream, _ := ioutil.ReadFile("alpr.conf")
	_ = json.Unmarshal(byteStream, &ntplog.Logs)

    file, err := os.Open("iptables.log")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
	
    scanner := bufio.NewScanner(file)
    ipWithLastVisit := make(map[string]string)

	lastBootTime, _ := getLastSystemBootTime()

    for scanner.Scan() {
        line := scanner.Text()
		list := strings.Split(line, "[")
		epch := strings.Trim(strings.Split(list[1], "]")[0], " ")
		epchfloat, err := strconv.ParseFloat(epch, 64)
		if (err == nil && strings.Contains(line, "[ntp]")){
			key := string(strings.Split(strings.Split(list[2], "SRC=")[1], " ")[0])
	
			lastEntry := time.Unix(lastBootTime.Unix()+int64(epchfloat), int64(100000*(epchfloat-float64(int64(epchfloat)))))
			if lastBootTime.Before(time.Now()) {
				lastEntryStr := lastEntry.String()
				lastEntryStr = strings.Split(lastEntryStr, " +")[0]
				ipWithLastVisit[key] = lastEntryStr
			}
		}
    }

	cnt := 0
    for _,logger := range ntplog.Logs {
		if timestamp, ok :=ipWithLastVisit[logger.Ip]; ok{
			ntplog.Logs[cnt].Time = timestamp
			cnt++
		}
	}
	ntplog.Logs = ntplog.Logs[:cnt]
	
	byteStream, _ = json.Marshal(ntplog)
	resp, err := http.Post("http://140.118.127.165:82/api/v1/ntp_logs", "application/json", bytes.NewBuffer(byteStream))
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(resp)
}
