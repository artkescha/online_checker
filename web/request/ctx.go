package request

import (
	"errors"
	"gitlab.com/artkescha/grader/online_checker/pkg/middlewares"
	"gitlab.com/artkescha/grader/online_checker/pkg/user"
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
