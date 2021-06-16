package validateiap

import (
	"context"
	"log"
	"net/http"

	"github.com/a1comms/gcp-iap-auth/jwt"
	"github.com/urfave/negroni"
)

type emailValFunc func(context.Context, string) (bool, error)

var (
	ValidateIAPMiddleware          negroni.HandlerFunc = GetValidateIAPMiddleware(emailNotEmpty)
	ValidateIAPAppEngineMiddleware negroni.HandlerFunc = GetValidateIAPAppEngineMiddleware(emailNotEmpty)
)

func GetValidateIAPMiddleware(emailVal emailValFunc) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if claims, err := jwt.RequestClaims(r, cfg); err != nil {
			log.Printf("ValidateIAP: Failed to validate request claims: %s", err)
		} else {
			if ok, err := emailVal(r.Context(), claims.Email); err != nil {
				log.Printf("ValidateIAP: Failed to call email validation function: %s", err)
			} else if ok {
				ctx, _ := setUserEmailToContext(r.Context(), claims.Email)
				ctx, _ = setGoogleClaimToContext(ctx, claims.Google)
				next(w, r.WithContext(ctx))
				return
			}
		}

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func GetValidateIAPAppEngineMiddleware(emailVal emailValFunc) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if val := r.Header.Get("X-AppEngine-Cron"); val != "" {
			next(w, r)
			return
		} else if val := r.Header.Get("X-AppEngine-QueueName"); val != "" {
			next(w, r)
			return
		} else if claims, err := jwt.RequestClaims(r, cfg); err == nil {
			if ok, err := emailVal(r.Context(), claims.Email); err != nil {
				log.Printf("ValidateIAP: Failed to call email validation function: %s", err)
			} else if ok {
				ctx, _ := setUserEmailToContext(r.Context(), claims.Email)
				ctx, _ = setGoogleClaimToContext(ctx, claims.Google)
				next(w, r.WithContext(ctx))
				return
			}
		} else {
			log.Printf("ValidateIAP: Failed to validate request claims: %s", err)
		}

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func emailNotEmpty(ctx context.Context, email string) (bool, error) {
	if email != "" {
		return true, nil
	}

	return false, nil
}
