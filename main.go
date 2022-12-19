package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

func main() {
	for {
		process()
		time.Sleep(time.Duration(getInterval()) * time.Second)
	}
}

func process() {
	err := mergeConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = realoadClashConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func mergeConfig() error {
	url := getUrl()
	if url == "" {
		panic(errors.New("url is empty"))
	}
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("get by url [%s] has error: %v", url, err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("get by url [%s] return: %d", url, resp.StatusCode)
	}
	defer resp.Body.Close()
	remote := viper.New()
	remote.SetConfigType("yaml")
	remote.ReadConfig(resp.Body)
	proxies := remote.Get("proxies")
	if proxies == nil {
		return fmt.Errorf("get proxies from url [%s] is empty", url)
	}
	proxyGroups := remote.Get("proxy-groups")
	rules := remote.Get("rules")
	local := viper.New()
	local.SetConfigFile(getPath())
	err = local.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read config has error: %v", err)
	}
	local.Set("proxies", proxies)
	if proxyGroups != nil {
		local.Set("proxy-groups", proxyGroups)
	}
	if rules != nil {
		local.Set("rules", rules)
	}
	err = local.WriteConfig()
	if err != nil {
		return fmt.Errorf("write config has error: %v", err)
	}
	log.Printf("merge config from url [%s] success", url)
	return nil
}

func realoadClashConfig() error {
	type configPut struct {
		Path string `json:"path"`
	}
	path := getClashConfigPathInClash()
	body := &configPut{
		Path: path,
	}
	bodyBytes, _ := json.Marshal(body)
	req, err := http.NewRequest("PUT", getClashUrl()+"/configs", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	secrect := getClashSecret()
	if secrect != "" {
		req.Header.Set("authorization", fmt.Sprintf("Bearer %s", secrect))
	}
	if err != nil {
		return errors.New("create request has error: " + err.Error())
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("request has error: " + err.Error())
	}
	if resp.StatusCode != 200 {
		return errors.New("request return: " + resp.Status)
	}
	return nil
}

func getClashConfigPathInClash() string {
	path := os.Getenv("CLASH_CONF_PATH_IN_CLASH")
	if path == "" {
		path = getPath()
	}
	return path
}

func getClashUrl() string {
	url := os.Getenv("CLASH_URL")
	if url == "" {
		url = "http://127.0.0.1:8080"
	}
	return url

}

func getClashSecret() string {
	secret := os.Getenv("CLASH_SECRET")
	return secret
}
func getInterval() int {
	interval := os.Getenv("INTERVAL")
	if interval == "" {
		interval = "3600"
	}
	result, err := strconv.Atoi(interval)
	if err != nil {
		log.Println("get interval has error: ", err, ", use default value 3600")
		result = 3600
	}
	if result <= 0 {
		log.Println("get interval is less than 0, use default value 3600")
		result = 3600
	}
	return result
}

func getPath() string {
	path := os.Getenv("CLASH_CONF_PATH")
	if path == "" {
		path = "/data/config.yaml"
	}
	return path
}

func getUrl() string {
	url := os.Getenv("CLASH_SUB_URL")
	return url
}
