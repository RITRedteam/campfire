package main

import "os/exec"
import "fmt"
import "log"
import "net/http"
import "io/ioutil"
import "bytes"

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
	cmd := exec.Command("whoami")
    out, err := cmd.CombinedOutput()
    if err != nil {
        log.Fatalf("cmd.Run() failed with %s\n", err)
        return "Err"
    }
    return string(out)

}

func send_data(rules string, host string){
	fmt.Printf("%s", host);
	fmt.Printf("%s", rules);
	url := "127.0.0.1:8880";

	bigStr := fmt.Sprintf("%s%s%s%s%s", "{\"hostname\":\"", host, "\", \"rules\":\"", rules, "\"}");
	fmt.Println("%s", bigStr);

    var jsonStr = []byte(bigStr);
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr));
    req.Header.Set("X-Custom-Header", "myvalue");
    req.Header.Set("Content-Type", "application/json");

    client := &http.Client{};
    resp, err := client.Do(req);
    if err != nil {
        panic(err);
    }
    defer resp.Body.Close();

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
}

func main(){
	rules := get_tables()
	host := get_hn()
	send_data(rules, host)
}
