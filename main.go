package main

import (
	_ "bufio"
	"fmt"
	"log"
	"flag"
	"html/template"
	_ "io"
	"os"
	"net/http"
	_ "strings"
    "unite"
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
var logserver  bool

func init() {
    flag.StringVar(&configpath, "f", "file.ini", "config file path")
	flag.BoolVar(&logserver, "l", false, "log server")
}

func main() {
	flag.Parse()
	if logserver {
		go unite.UDPLogServer(1111, "127.0.0.1")
	}

	data := struct {
		Items []string
	}{
		Items: []string{},
	}

	c := unite.ParseConfig(configpath)

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
		path = fmt.Sprintf("/%s", k)
		http.Handle(path, unite.FileSummaryServer(items[k], k))

	}

	t, err := template.New("webpage").Parse(tpl)
	if err != nil {
		fmt.Print(err)
		return
	}
	
	fmt.Printf("reload config: kill -SIGHUP %d\n", os.Getpid())

	handler := func(w http.ResponseWriter, r *http.Request) {
		_ = t.Execute(w, data)
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
