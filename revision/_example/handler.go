package main

import (
	"net/http"

	"github.com/motain/samodelkin/revision"
)

func main() {
	// Revision file created during build process
	// check Makefile example
	r := revision.AppRevision()
	if err := r.Load(); err != nil {
		panic(err)
	}

	http.Handle("/_healthcheck_", NewHealthcheckHandler(r))
	http.ListenAndServe(":8080", nil)
}

type HealthcheckHandler struct {
	r revision.AppRevision
}

func NewHealthcheckHandler(r AppRevision) *HealthcheckHandler {
	return &HealthcheckHandler{r}
}

func (h *HealthcheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// write a current commit hash in response
	w.Write(h.r.Message())
}
