package main

import (
    "fmt"
    "os/exec"
)

func main() {
    cmd := "cat /var/log/iptables.log"
    out, err := exec.Command("bash", "-c", cmd).Output()
    if err == nil {
        fmt.Print(out)
    }
}