package session

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/artkescha/checker/online_checker/pkg/user"
	"time"
)

const ExpireTimeDay = 7

//TODO потом удалить!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
const SecretKey = ""

type Claims struct {
	User      user.User `json:"user"`
	SessionID string    `json:"session"`
	*jwt.StandardClaims
}

func createJWT(sessionID string, user user.User) (string, error) {
	d := time.Duration(time.Hour * 24 * ExpireTimeDay)
	expireToken := time.Now().Add(d).Unix()
	claims := Claims{
		user,
		sessionID,
		&jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expireToken,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//SecretKey, exists := os.LookupEnv("SECRET_KEY")
	//if !exists {
	//	return "", fmt.Errorf("SECRET_KEY is required")
	//}
	signedToken, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", fmt.Errorf("signed token failed: %s", err)
	}
	return signedToken, nil
}
