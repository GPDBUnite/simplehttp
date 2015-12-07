package main

import (
	_ "bufio"
	"fmt"
	"html/template"
	_ "io"
	"net/http"
	"os"
	_ "strings"
    "unite"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func main() {
	data := struct {
		Items []string
	}{
		Items: []string{},
	}

	filename := `file.ini`
	c := unite.ParseConfig(filename)
    // fmt.Println(c)
	items := c["static"]
	for k := range items {
		ok, _ := exists(items[k])
		if !ok {
			fmt.Printf("skip %s, not exist\n", k)
			continue
		}
		data.Items = append(data.Items, k)
		path := fmt.Sprintf("/%s/", k)
		http.Handle(path, http.StripPrefix(path, http.FileServer(http.Dir(items[k]))))
		//fmt.Printf("%s=>%s\n", k, items[k])
	}

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
	t, err := template.New("webpage").Parse(tpl)
	if err != nil {
		fmt.Print(err)
		return
	}
	
	fmt.Printf("reload config: kill -SIGHUP %d\n", os.Getpid())
	// log.Println("Listening...")
	handler := func(w http.ResponseWriter, r *http.Request) {
		_ = t.Execute(w, data)
	}
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
