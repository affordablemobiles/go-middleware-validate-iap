package validateiap

import (
	"fmt"
	"net/http"
	"strings"
)

func GetUserEmail(r *http.Request) (string, error) {
	email := r.Header.Get("X-Goog-Authenticated-User-Email")
	if email == "" {
		return "", fmt.Errorf("Authenticated email is blank")
	}

	return strings.Replace(email, "accounts.google.com:", "", 1), nil
}
