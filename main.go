package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	url    = "https://en.wikipedia.org/w/api.php?origin=*"
	params = "action=query&list=recentchanges&rcprop=title|ids|sizes|flags|user&rclimit=3&format=json"
)

func main() {
	request, err := http.Get(url + "&" + params)
	if err != nil {
		panic(err)
	}

	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}
