package main

import (
	"encoding/xml"
	"encoding/json"
	"io/ioutil"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

// Feed is an xml RSS feed. We only care about titles here
type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Titles  []string `xml:"channel>item>title"`
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
	numWords := flag.Int("words", 25, "maximum number of words to print")
	prefixLen := flag.Int("prefix", 2, "prefix length in words")

	flag.Parse()                     // Parse command-line flags.
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator.

	titles, err := Load()

	if err != nil {
		log.Fatal(err)
	}

	c := NewChain(*prefixLen) // Initialize a new Chain.

	for title, _ := range titles {
		c.Build(title)
	}

	text := c.Generate(*numWords) // Generate text.

	_, ok := titles[text]

	fmt.Println(text)

	if !ok {
		fmt.Println("NEW HEADLINE")
	} 
}
