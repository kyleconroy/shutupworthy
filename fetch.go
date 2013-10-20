
package main

import (
	"encoding/xml"
	"encoding/json"
	"net/http"
	"log"
	"io/ioutil"
)

// Feed is an xml RSS feed. We only care about titles here
type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Titles  []string `xml:"channel>item>title"`
}

// Return the current titles to Upworthy blog posts
func FetchTitles() ([]string, error) {
	v := Feed{}

	resp, err := http.Get("http://feeds.feedburner.com/upworthy")

	if err != nil {
		return []string{}, err
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return []string{}, err
	}

	defer resp.Body.Close()

	err = xml.Unmarshal([]byte(data), &v)

	if err != nil {
		return []string{}, err
	}

	return v.Titles, nil
}

func Load() (map[string]string, error) {
	var archive map[string]string

	blob, err := ioutil.ReadFile("titles.json")

	if err != nil {
		return map[string]string{}, err 
	}

	err = json.Unmarshal(blob, &archive)

	if err != nil {
		return map[string]string{}, err 
	}

	return archive, nil
}

func Save(archive map[string]string) error {
	blob, err := json.Marshal(archive)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile("titles.json", blob, 0644)

	if err != nil {
		return err
	}

	return nil
}


func main() {
	//Load current titles
	archive, err := Load()

	if err != nil {
		log.Fatal(err)
	}

	//Fetch new ones
	titles, err := FetchTitles()

	if err != nil {
		log.Fatal(err)
	}

	//Merge
	for _, title := range titles {
		archive[title] = ""
	}

	err = Save(archive)

	if err != nil {
		log.Fatal(err)
	}
}
