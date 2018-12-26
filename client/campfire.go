// This go file will collect the host's firewall rules ship them back
// to a defined webserver, along with the hostname
// disclaimer: this doesn't work yet
// @author: degenerat3
package main

import "os/exec"
import "fmt"
import "log"
import "net/http"
import "io/ioutil"
import "bytes"
import "encoding/json"

 var serv = "127.0.0.1:5000/api/rule_send"
// var loop_time = 30

func get_tables() string{
	cmd := exec.Command("iptables", "-L")
    out, err := cmd.CombinedOutput()
    if err != nil {
        log.Fatalf("cmd.Run() failed with %s\n", err)
		return "Err"
    }
	return string(out)
}

func get_hn() string{
	cmd := exec.Command("hostname")
    out, err := cmd.CombinedOutput()
    if err != nil {
        log.Fatalf("cmd.Run() failed with %s\n", err)
        return "Err"
    }
    return string(out)

}

func send_data(rules string, host string){
	url1 := "http://" + serv;
	fmt.Printf("url: %s\n", url1)
	bigStr := fmt.Sprintf("%s%s%s%s%s", "{\"hostname\":\"", host, "\", \"rules\":\"", rules, "\"}");
	fmt.Println("%s", bigStr);

    jsonData := map[string]string{"hostname": host, "rules": rules}
	jsonValue, _ := json.Marshal(jsonData)
	resp, err := http.Post(url1, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("Req failed: %s\n", err)
	} else{
		data, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(data))
	}
}

func main(){
	rules := get_tables()
	host := get_hn()
	send_data(rules, host)
}
