package unite

import (
	_ "fmt"
	"net/http"
	_ "path/filepath"
)

type mdhandler struct {
	root   string
}

func (f *mdhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
     w.Write([]byte("hello"))
}

func MarkDownHandler(root string) http.Handler {
	return &mdhandler{root}
}
