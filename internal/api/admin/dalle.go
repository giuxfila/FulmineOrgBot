package admin

import (
	"github.com/giuxfila/FulmineOrgBot/internal/dalle"
	"net/http"
)

func (s Service) DisableDalle(w http.ResponseWriter, r *http.Request) {
	dalle.Enabled = false
}

func (s Service) EnableDalle(w http.ResponseWriter, r *http.Request) {
	dalle.Enabled = true
}
