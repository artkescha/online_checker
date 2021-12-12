package middlewares

import (
	"github.com/artkescha/checker/online_checker/pkg/session"
	"github.com/artkescha/checker/online_checker/web/response"
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
		//TODO это костыль! Нужно поправить!
		if token == "" {
			token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7ImlkIjoyLCJ1c2VybmFtZSI6ImFkbWluIn0sInNlc3Npb24iOiI4ZGRiNTEzMy0yNTNhLTQ1MTgtYWM4Zi1jYzFkOWRmOGRlY2IiLCJleHAiOjE2Mzk4MzY4MzksImlhdCI6MTYzOTIzMjAzOX0.eblMZww8isWEl2JK2Ject9kMWeKsRjW8nXvHMm9LIqE"
		}
		session, err := manager.GetSession(token)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		ctx := context.WithValue(r.Context(), ContextSession, session.User)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
