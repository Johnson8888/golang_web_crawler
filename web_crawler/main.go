package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	sendRequest()
}

func sendRequest() {
	webUrl := "https://bet.hkjc.com/football/getJSON.aspx?jsontype=schedule.aspx"
	resp, err := http.Get(webUrl)
	if err != nil {
		panic("request failed with error " + err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
