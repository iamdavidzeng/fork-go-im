package wordsfilter

import (
	"bufio"
	"log"
	"os"

	"github.com/syyongx/go-wordsfilter"
)

var (
	Wf      *wordsfilter.WordsFilter
	samples []string
	root    map[string]*wordsfilter.Node
)

func SetTexts() {
	f, err := os.Open("sample.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := scanner.Text()
		samples = append(samples, s)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	Wf = wordsfilter.New()
	root = Wf.Generate(samples)
}

func MsgFilter(val string) bool {
	return Wf.Contains(val, root)
}
