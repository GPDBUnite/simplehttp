package main

import (
	"github.com/lij55/unite"
	_ "bufio"
	"flag"
	"fmt"
	"html/template"
	_ "io"
	"log"
	"net/http"
	"os"
	"strconv"
	_ "strings"
)

const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Static files</title>
	</head>
	<body>
		{{range .Items}}
			<div><a href="/{{ . }}">{{ . }}</a></div>
		{{end}}
	</body>
</html>`

var configpath string
var logserver bool
var s3list bool
var mdlist bool


func init() {
	flag.StringVar(&configpath, "f", "file.ini", "config file path")
	flag.BoolVar(&logserver, "l", false, "log server")
  s3list = true
  mdlist = true
}

func main() {
	flag.Parse()

	data := struct {
		Items []string
	}{
		Items: []string{},
	}

	c := unite.ParseConfig(configpath)

	if logserver {
		val := c.GetConfig("log", "port", "1111")
		loghost := c.GetConfig("log", "host", "127.0.0.1")
		port, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		//fmt.Println(loghost)
		//fmt.Println(val)
		go unite.UDPLogServer(port, loghost)
	}

	items := c["static"]

	for k := range items {
		ok, _ := unite.FileExists(items[k])
		if !ok {
			fmt.Printf("skip %s, not exist\n", k)
			continue
		}
		data.Items = append(data.Items, k)
		path := fmt.Sprintf("/%s/", k)
		http.Handle(path, http.StripPrefix(path, http.FileServer(http.Dir(items[k]))))
    if s3list {
		   path = fmt.Sprintf("/%s", k)
		   http.Handle(path, unite.FileSummaryServer(items[k], k))
    }
	}

  items = c["markdown"]
  for k := range items {
      ok, _ := unite.FileExists(items[k])
      if !ok {
         fmt.Printf("skip %s, not exist\n",k)
         continue
      }
      http.Handle(path, http.StripPrefix(path, unite.MarkDownHandler(items[k])))
  }

	t, err := template.New("webpage").Parse(tpl)
	if err != nil {
		fmt.Print(err)
		return
	}

	//fmt.Printf("reload config: kill -SIGHUP %d\n", os.Getpid())

	handler := func(w http.ResponseWriter, r *http.Request) {
		_ = t.Execute(w, data)
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
