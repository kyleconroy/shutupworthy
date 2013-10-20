package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
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

// Prefix is a Markov chain prefix of one or more words.
type Prefix []string

// String returns the Prefix as a string (for use as a map key).
func (p Prefix) String() string {
	return strings.Join(p, " ")
}

// Shift removes the first word from the Prefix and appends the given word.
func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

// Chain contains a map ("chain") of prefixes to a list of suffixes.
// A prefix is a string of prefixLen words joined with spaces.
// A suffix is a single word. A prefix can have multiple suffixes.
type Chain struct {
	chain     map[string][]string
	prefixLen int
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func NewChain(prefixLen int) *Chain {
	return &Chain{make(map[string][]string), prefixLen}
}

// Build reads text from the provided Reader and
// parses it into prefixes and suffixes that are stored in Chain.
func (c *Chain) Build(line string) {
	p := make(Prefix, c.prefixLen)

	for _, s := range strings.Split(line, " ") {
		key := p.String()
		c.chain[key] = append(c.chain[key], s)
		p.Shift(s)
	}
}

// Generate returns a string of at most n words generated from Chain.
func (c *Chain) Generate(n int) string {
	p := make(Prefix, c.prefixLen)
	var words []string
	for i := 0; i < n; i++ {
		choices := c.chain[p.String()]
		if len(choices) == 0 {
			break
		}
		next := choices[rand.Intn(len(choices))]
		words = append(words, next)
		p.Shift(next)
	}
	return strings.Join(words, " ")
}

func main() {
	// Register command-line flags.
	numWords := flag.Int("words", 20, "maximum number of words to print")
	prefixLen := flag.Int("prefix", 2, "prefix length in words")

	flag.Parse()                     // Parse command-line flags.
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator.

	titles, err := FetchTitles()

	if err != nil {
		log.Fatal(err)
	}

	c := NewChain(*prefixLen) // Initialize a new Chain.

	for _, title := range titles {
		c.Build(title)
	}

	text := c.Generate(*numWords) // Generate text.
	fmt.Printf("Chain: %v\n", c)
	fmt.Println(text)             // Write text to standard output.
}
