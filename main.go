package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type storyArc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var chapters map[string]storyArc

const tmpl = `
    <h1>{{.Title}}</h1>
	<p>{{.Story}}</p>
	<ul>
    {{range .Options}}
	<li><a href="/{{.Arc}}">{{.Text}}</a></li>
    {{end}}
	</ul>`

type MyHandler struct{}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("webpage").Parse(tmpl)
	check(err)

	path := r.URL.Path
	if path == "/" {
		t.Execute(w, chapters["intro"])
	} else {
		t.Execute(w, chapters[strings.Trim(path, "/")])
	}
}

func (h *MyHandler) ServeCLI(path string) {
	chapter := chapters[path]

	fmt.Println(chapter.Title)

	for _, item := range chapter.Story {
		fmt.Println(item)
	}

	fmt.Println("")

	if len(chapter.Options) == 0 {
		return
	}

	opts := make(map[int]string)

	for key, item := range chapter.Options {
		opts[key+1] = item.Arc
		fmt.Printf("%d: %s\n", key+1, item.Text)
	}

	reader := bufio.NewReader(os.Stdin)
	input, _, _ := reader.ReadLine()
	selected, _ := strconv.Atoi(string(input))

	for {
		if opts[selected] == "" {
			fmt.Println("Wrong input")
			input, _, _ = reader.ReadLine()
			selected, _ = strconv.Atoi(string(input))
		} else {
			break
		}
	}

	h.ServeCLI(opts[selected])
}

func main() {

	cli := flag.Bool("cli", false, "use cli instead of web version")
	flag.Parse()

	data, err := os.ReadFile("gopher.json")
	check(err)

	err = json.Unmarshal(data, &chapters)
	check(err)

	if *cli {
		Handler := MyHandler{}
		Handler.ServeCLI("intro")
	} else {
		mux := http.NewServeMux()
		mux.Handle("/", &MyHandler{})

		println("Serving at http://localhost:8080")
		http.ListenAndServe(":8080", mux)
	}
}
