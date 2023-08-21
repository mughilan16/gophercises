package main

import (
	"cyoa"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

func main() {
	port := flag.Int("port", 3000, "the port to start the CYOA web application on")
	filename := flag.String("file", "gopher.json", "the JSON file with the CYOA story")
	flag.Parse()
	fmt.Printf("Using the story in %s.\n", *filename)

	f, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}
	story, err := cyoa.JsonStory(f)
	if err != nil {
		panic(err)
	}
	tpl := template.Must(template.New("").Parse(tmpl))
	h := cyoa.NewHandler(story,
		cyoa.WithPathFunc(pathFn),
		cyoa.WithTemplate(tpl),
	)
	mux := http.NewServeMux()
	mux.Handle("/story/", h)
	fmt.Printf("Starting the server at: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}

func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}
	return path[len("/story/"):]
}

var tmpl = `
<!DOCTYPE html>
<html>
  <head>
		<link rel="preconnect" href="https://fonts.googleapis.com">
		<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
		<link href="https://fonts.googleapis.com/css2?family=Poppins:wght@300;500&display=swap" rel="stylesheet">
    <meta charset="utf-8">
    <title>Choose Your Own Adventure</title>
		<style>
			body {
				font-family: 'Poppins', sans-serif;
			}
			h1 {
				text-align: center;
				position: relative;
			}
			.page {
				width: 80%;
				max-width: 500px;
				margin: auto;
				margin-top: 40px;
				margin-bottom: 40px;
				padding: 80px;
				background: #FFFCF6;
				border: 1px solid #eee;
				box-shadow: 0 8px 6px -6px #777;
			}
			ul {
				border-top: 1px dotted #ccc;
				padding: 10px 0 0 0;
				-webkit-padding-start: 0;
			}
			li {
				padding-top: 10px;
			}
			a,
			a:visited {
				text-decoration: none;
				color: #6295b5;
			}
			a:active,
			a:hover {
				color: #7792a2;
			}
			p {
				text-indent: 1em;
			}
		</style>
  </head>
  <body>
		<section class="page">
    <h1>{{.Title}}</h1>
    {{range .Paragraphs}}
    <p>{{.}}</p>
    {{end}}
    <ul>
    {{range .Options}}
      <li><a href="/story/{{.Chapter}}">{{.Text}}</a></li>
    {{end}}
    </ul>
		</section>
  </body>
</html>
`
