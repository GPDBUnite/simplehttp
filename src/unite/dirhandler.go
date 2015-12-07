package unite

import (
	"net/http"
	"fmt"
	"path/filepath"
)


type FileSummary struct {
	root string
}

func (f *FileSummary) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	rootfs := http.Dir(f.root)
	fs, err := rootfs.Open(".")
	if err != nil {
		// error
		return
	}
	for {
		dirs, err := fs.Readdir(100)
		if err != nil || len(dirs) == 0 {
			break
		}
		for _, d := range dirs {
			name := d.Name()
			
			if d.IsDir() {
				continue
			}
			fullpath := filepath.Join(f.root,name)
			size, err := FileSize(fullpath)
			if err == nil {
				fmt.Fprintf(w, "%s,%d\n",name, size)
			}
		}
	}
	
}

func FileSummaryServer(root string) http.Handler {
	return &FileSummary{root}
}
