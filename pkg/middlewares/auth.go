package middlewares

import (
	"gitlab.com/artkescha/grader/online_checker/pkg/session"
	"gitlab.com/artkescha/grader/online_checker/web/response"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

//authorization

type ContextKey string

const ContextSession ContextKey = "user"

func Authorization(manager session.Manager, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		session, err := manager.GetSession(token)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		ctx := context.WithValue(r.Context(), ContextSession, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
