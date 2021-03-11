package request

import (
	"errors"
	"github.com/artkescha/checker/online_checker/pkg/middlewares"
	"github.com/artkescha/checker/online_checker/pkg/user"
	"net/http"
)

func ExtractContext(r *http.Request) (user.User, error) {
	ctx := r.Context().Value(middlewares.ContextSession)
	user, ok := ctx.(user.User)
	if !ok {
		return user, errors.New("request context failed")
	}
	return user, nil
}
