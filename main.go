package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

const (
	url     = "https://fr.wikipedia.org/w/api.php?origin=*"
	params  = "action=query&list=recentchanges&rcprop=title|ids|sizes|flags|user|timestamp&rclimit=100&format=json"
	ipRegex = `((^\s*((([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]))\s*$)|(^\s*((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(%.+)?\s*$))`
)

type response struct {
	Query query `json:"query"`
}

type query struct {
	RecentChanges []change `json:"recentchanges"`
}

type change struct {
	Type      string `json:"type"`
	Title     string `json:"title"`
	User      string `json:"user"`
	OldLength string `json:"oldlen"`
	NewLength string `json:"newlen"`
	Timestamp string `json:"timestamp"`
}

func main() {
	for {
		for _, change := range getRecentChanges() {
			if matchIP(change.User) {
				fmt.Println(change)
			}
		}

		time.Sleep(30 * time.Second)
	}
}

func matchIP(username string) bool {
	regex := regexp.MustCompile(ipRegex)

	if !regex.MatchString(username) {
		return false
	}

	for i := 17; i <= 30; i++ {
		if username != fmt.Sprintf("195.24.201.%d", i) {
			continue
		}

		return true
	}

	return false
}

func getRecentChanges() []change {
	request, err := http.Get(url + "&" + params)
	if err != nil {
		panic(err)
	}

	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}

	var data response
	json.Unmarshal(body, &data)

	return data.Query.RecentChanges
}
