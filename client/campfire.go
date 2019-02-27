// This go file will collect the host's firewall rules and other networking info, ship them back to a defined webserver,
// along with the hostname/ip
// disclaimer: this barely works
// @author: degenerat3
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var serv = "127.0.0.1:5000" //IP of flask serv
var loopTime = 500          //sleep time in seconds

// return output of "iptables -L" as one large string
func getTables() string {
	cmd := exec.Command("/bin/bash", "-c", "echo \"Filter Table\"; iptables -t filter -L; echo; echo; echo \"NAT Table\"; iptables -t nat -L; echo; echo; echo \"Mangle Table\"; iptables -t mangle -L; echo; echo; echo \"Raw Table\"; iptables -t raw -L;")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return "Err"
	}
	return string(out)
}

// return output of "iptables -L" as one large string
func getHosts() string {
	cmd := exec.Command("/bin/bash", "-c", "cat /etc/hosts")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return "Err"
	}
	return string(out)
}

// return output of "iptables -L" as one large string
func getRoutes() string {
	cmd := exec.Command("/bin/bash", "-c", "ip route")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return "Err"
	}
	return string(out)
}

// return output of "iptables -L" as one large string
func getArp() string {
	cmd := exec.Command("/bin/bash", "-c", "arp -a")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return "Err"
	}
	return string(out)
}

// return hostname as string
func getHn() string {
	hstr, _ := os.Hostname()
	return hstr

}

func getIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	ad := conn.LocalAddr().(*net.UDPAddr)
	ipStr := ad.IP.String()
	ipStr = strings.Replace(ipStr, ".", "-", -1)
	return ipStr
}

// post strings to flask server
func sendData(rules string, hosts string, routes string, arp string, host string, ip string) {
	url1 := "http://" + serv + "/api/rule_send" // turn ip into valid url
	jsonData := map[string]string{"rules": rules, "etchosts": hosts, "routes": routes, "arp": arp, "hostname": host, "ip": ip}
	jsonValue, _ := json.Marshal(jsonData)
	insRule := exec.Command("iptables", "-I", "FILTER", "1", "-j", "ACCEPT") //temporarily allow so we can send data
	insRule.Run()
	_, err := http.Post(url1, "application/json", bytes.NewBuffer(jsonValue))
	dropRule := exec.Command("iptables", "-D", "FILTER", "1")
	dropRule.Run()
	if err != nil {
		fmt.Printf("Req failed: %s\n", err)
		return
	} else {
		// block below for debug
		// data, _ := ioutil.ReadAll(resp.Body)
		// fmt.Println(string(data))
		return
	}
}

// fetch data then send it
func run() {
	rules := getTables()
	hosts := getHosts()
	routes := getRoutes()
	arp := getArp()
	host := getHn()
	ip := getIP()
	sendData(rules, hosts, routes, arp, host, ip)
}

func main() {
	loopArg := os.Args[1]
	if loopArg == "-s" { // if "-s" is an arg, run once, otherwise loop
		run()
	} else {
		for { // send data to webserver ever X seconds, until termination
			run()
			t := time.Duration(loopTime * 1000)
			time.Sleep(t * time.Millisecond)
		}
	}
}
