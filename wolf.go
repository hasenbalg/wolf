package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"gopkg.in/yaml.v2"
)

/////////////////
// / model
/////////////////

type Config struct {
	ThisApplicationPort string `yaml:"thisApplicationPort"`
	LoginEndoint        string `yaml:"loginEndoint"`
	MacAddress          string `yaml:"macAddress"`
	BroadcastAddress    string `yaml:"broadcastAddress"`
}

/////////////////
// / config file
/////////////////

func readConfigFile() {
	/// which paths will be searched for a config file
	configPaths := []string{
		"/etc/wolf/config.yaml", "~/.config/wolf/config.yaml",
		"config.yaml"}

	var c Config
	for _, configPath := range configPaths {
		// try each config file
		yamlFileData, err := os.ReadFile(configPath)
		if err != nil {
			continue
		}

		// pare to config model
		err = yaml.Unmarshal(yamlFileData, &c)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
	}
	// could a config file be parsed
	// if  == nil {
	// 	fmt.Printf("cannot find config in: %s\n", configPaths)
	// 	os.Exit(1)
	// }
	config = &Config{
		ThisApplicationPort: c.ThisApplicationPort,
		LoginEndoint:        c.LoginEndoint,
		MacAddress:          c.MacAddress,
		BroadcastAddress:    c.BroadcastAddress,
	}

}

/////////////////
/// controllers
/////////////////

func getRoot(w http.ResponseWriter, r *http.Request) {
	// t := template.New("templates/index.html")   // Create a template.
	// t, _ = t.ParseFiles("templates/index.html") // Parse template file.

	t := template.Must(template.ParseFiles("templates/index.html"))

	// user := GetUser() // Get current user infomration.
	// t.Execute(w, user)  //
	t.Execute(w, nil) //
	fmt.Println(t)
	fmt.Printf("got / request\n")
	// io.WriteString(w, "This is my website!\n")
}

func getWake(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("/usr/bin/wakeonlan", "-i", config.BroadcastAddress, config.MacAddress)
	fmt.Println("IP-", config.BroadcastAddress)

	output, err := cmd.Output()

	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}

	fmt.Println(string(output))
	io.WriteString(w, "wake\n")
}

func getPing(w http.ResponseWriter, r *http.Request) {
	/// make http request with timeout
	fmt.Println(config.LoginEndoint)

	customTransport := http.DefaultTransport.(*http.Transport)

	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := http.Client{Transport: customTransport, Timeout: 5 * time.Second}

	resp, err := client.Get(config.LoginEndoint)
	if err != nil {
		fmt.Println(err)
		fmt.Println("false1")
		io.WriteString(w, "false")
		return
	}
	fmt.Println(resp.StatusCode)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			// fmt.Println(err)
			fmt.Println("false")
			io.WriteString(w, "false")
			return
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
		fmt.Println("true")
		io.WriteString(w, "true")
	} else {

		fmt.Println("false")
		io.WriteString(w, "false")
	}

}

// / global instance of config
var config *Config

func main() {
	// read config file
	readConfigFile()

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/wake", getWake)
	http.HandleFunc("/ping", getPing)

	err := http.ListenAndServe(fmt.Sprintf(":%s", config.ThisApplicationPort), nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
