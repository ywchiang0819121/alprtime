package main

import (
    "fmt"
    "os/exec"
)

func main() {
    cmd := "cat /var/log/iptables.log"
    out, err := exec.Command("bash", "-c", cmd).Output()
    if err == nil {
		// fmt.Println(string(out))
		var host string
		var timestmp string
		var etw string
		while(fmt.Sscanf(string(out), "%s gamelab-MS-7A68 %s SRC=%s DST%s", &timestmp, &etw, &host, &etw)){
			fmt.Println(timestmp, host)
		}
    }
}