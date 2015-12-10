package unite

import (
	"fmt"
	"net/http"
	"path/filepath"
)

type FileSummary struct {
	root   string
	bucket string
}

func (f *FileSummary) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	rootfs := http.Dir(f.root)
	fs, err := rootfs.Open(".")
	if err != nil {
		// error
		return
	}
	fmt.Fprintf(w, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	fmt.Fprintf(w, "<ListBucketResult xmlns=\"http://s3.amazonaws.com/doc/2006-03-01/\">")
	fmt.Fprintf(w, "<Name>%s</Name>", f.bucket)
	fmt.Fprintf(w, "<Prefix></Prefix>")
	//fmt.Fprintf(w, "")
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
			fullpath := filepath.Join(f.root, name)
			size, err := FileSize(fullpath)
			if err == nil {
				fmt.Fprintf(w, "<Contents>")
				fmt.Fprintf(w, "<Key>%s</Key>", name)
				fmt.Fprintf(w, "<Size>%d</Size>", size)
				fmt.Fprintf(w, "</Contents>")
			}
		}
	}
	fmt.Fprintf(w, "</ListBucketResult>")

}

func FileSummaryServer(root string, bucket string) http.Handler {
	return &FileSummary{root, bucket}
}
