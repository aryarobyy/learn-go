package controller

import (
	"net/http"

	"github.com/aryarobyy/learn-go/config"
)

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := config.AppConfig.GoogleLoginConfig.AuthCodeURL("randomstate")
	http.Redirect(w, r, url, http.StatusSeeOther)
}
