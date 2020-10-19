package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Option struct {
	ArcName string `json:"arc"`
	Text    string `json:"text"`
}

type Arc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

const tPath = "resources/adventure.html"
const storyPath = "resources/gopher.json"

func getStoryHandler(arcs map[string]Arc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		params, ok := r.URL.Query()["next"]
		var arc Arc
		if !ok {
			arc = arcs["intro"]
		} else {
			key := params[0]
			arc, ok = arcs[key]
			if !ok {
				w.WriteHeader(404)
				return
			}
		}

		t, err := template.ParseFiles(tPath)
		if err != nil {
			log.Fatalf("could not parse template %s: %v", tPath, err)
		}
		err = t.Execute(w, arc)
		if err != nil {
			log.Fatal("could not insert data into template")
		}
	}
}

func parseStory(storyPath string) map[string]Arc {
	bytes, err := ioutil.ReadFile(storyPath)
	if err != nil {
		log.Fatalf("could not open %s", storyPath)
	}

	var preParsed map[string]json.RawMessage

	err = json.Unmarshal(bytes, &preParsed)
	if err != nil {
		log.Fatal("could not unmarshall bytes")
	}

	arcs := map[string]Arc{}

	for arcName, bytes := range preParsed {
		arc := Arc{}
		err := json.Unmarshal(bytes, &arc)
		if err != nil {
			log.Fatal("could not unmarshal arc")
		}
		arcs[arcName] = arc
	}

	return arcs
}

func main() {
	arcs := parseStory(storyPath)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./resources/static"))))
	http.HandleFunc("/", getStoryHandler(arcs))
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
