package main

import (
	"flag"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/russross/blackfriday"
)

var (
	useEmbeddedCss = flag.Bool("embedded-css", false, "")
	author         = flag.String("author", "", "author")
	title          = flag.String("title", "", "document title")
	writtenAt      = flag.String("at", "", "datetime, format: YYYY-mm-dd HH:MM:SS, default: now")
	inFile         = flag.String("i", "", "input markdown file, default: stdin")
	outFile        = flag.String("o", "", "output file, default: stdout")
)

type Doc struct {
	CSS         bool
	MdContent   string
	HtmlContent string
	Datetime    string
	Author      string
	Title       string
}

func main() {
	flag.Parse()

	html, err := Asset("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	funcMap := template.FuncMap{
		"attr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s)
		},
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	tmpl, err := template.New("tmpl").Funcs(funcMap).Parse(string(html))
	if err != nil {
		log.Fatal(err)
	}

	var md []byte
	if len(*inFile) == 0 {
		md, err = ioutil.ReadAll(os.Stdin)
	} else {
		md, err = ioutil.ReadFile(*inFile)
	}
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	err = tmpl.Execute(os.Stdout, Doc{
		CSS:         *useEmbeddedCss,
		MdContent:   string(md),
		HtmlContent: string(blackfriday.MarkdownCommon(md)),
		Datetime:    time.Now().Format("2006-01-02 15:04:05"),
		Author:      *author,
		Title:       *title,
	})

	if err != nil {
		log.Fatal(err)
	}

}
