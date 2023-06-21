package validateiap

import (
	"net/http"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/_gcp_iap/clear_login_cookie", 302)
}
