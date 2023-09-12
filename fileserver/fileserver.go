package fileserver

import (
	"io"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"path"
	"text/template"
)

func WithFallback(root fs.FS, basePath, fallbackPath string, templateData any) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := path.Join(basePath, path.Clean(r.URL.Path)[1:])

		info, err := fs.Stat(root, p)
		if err != nil || info.IsDir() {
			t, err := template.ParseFS(root, path.Join(basePath, fallbackPath))
			if err != nil {
				log.Print(err)
				return
			}

			w.Header().Add("Content-Type", "text/html")
			err = t.Execute(w, templateData)
			if err != nil {
				log.Print(err)
				return
			}
			return
		}

		f, err := root.Open(p)
		if err != nil {
			log.Print(err)
			return
		}

		w.Header().Add("Content-Type", mime.TypeByExtension(path.Ext(p)))

		_, err = io.Copy(w, f)
		if err != nil {
			log.Print(err)
			return
		}
	})
}
