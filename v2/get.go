package validateiap

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/a1comms/gcp-iap-auth/jwt"
)

type ctx_key int

const (
	ctx_key_user_email   ctx_key = iota
	ctx_key_google_claim ctx_key = iota
)

func GetUserEmail(r *http.Request) (string, error) {
	email := r.Header.Get("X-Goog-Authenticated-User-Email")
	if email == "" {
		return "", fmt.Errorf("Authenticated email is blank")
	}

	return strings.Replace(email, "accounts.google.com:", "", 1), nil
}

func GetUserEmailFromContext(ctx context.Context) (string, error) {
	if email, ok := ctx.Value(ctx_key_user_email).(string); ok {
		if email == "" {
			return "", fmt.Errorf("Authenticated email is blank")
		}

		return strings.Replace(email, "accounts.google.com:", "", 1), nil
	} else {
		return "", fmt.Errorf("Failed to fetch a valid email from the context")
	}
}

func GetGoogleClaimFromContext(ctx context.Context) (jwt.GoogleClaim, error) {
	if claim, ok := ctx.Value(ctx_key_google_claim).(jwt.GoogleClaim); ok {
		return claim, nil
	} else {
		return jwt.GoogleClaim{}, fmt.Errorf("Failed to fetch a valid email from the context")
	}
}

func setUserEmailToContext(ctx context.Context, email string) (context.Context, error) {
	return context.WithValue(
		ctx,
		ctx_key_user_email,
		email,
	), nil
}

func setGoogleClaimToContext(ctx context.Context, claim jwt.GoogleClaim) (context.Context, error) {
	return context.WithValue(
		ctx,
		ctx_key_google_claim,
		claim,
	), nil
}
