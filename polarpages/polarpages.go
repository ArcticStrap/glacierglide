package polarpages

import (
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Setup(rt *chi.Mux,addr string) {
  rt.Get("/wiki/{title}",func(w http.ResponseWriter,r *http.Request) {
    titleParam := chi.URLParam(r,"title")

    // Get page content
    res, err := http.Get(addr + "/api/wiki/" + titleParam)
    if err != nil {
      http.Error(w,err.Error(),http.StatusInternalServerError)
      return
    }
    defer res.Body.Close()

    // Read response
    content, err := io.ReadAll(res.Body)
    if err != nil {
      http.Error(w,err.Error(),http.StatusInternalServerError)
      return
    }

    // Parse template
    tmpl := template.Must(template.ParseFiles("polarpages/templates/index.html"))
    tmpl.Execute(w,struct{Title string 
      Content template.HTML}{
        Title: titleParam,
        Content: template.HTML(content),
      })
  })

  log.Println("PolarPages initalized")
}
