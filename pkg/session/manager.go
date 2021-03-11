package session

import (
	"fmt"
	"github.com/artkescha/checker/online_checker/pkg/user"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"reflect"

	"strconv"
)

const sessionLifeDay = 7

type Session struct {
	ID   string
	User user.User
}

//go:generate mockgen -destination=./session_mock.go -package=session . Manager

type Manager interface {
	CreateSession(user user.User) (string, error)
	GetSession(token string) (*Session, error)
	DestroySession(sessionId string) error
}

type SessionManager struct {
	client *memcache.Client
}

func NewManager(client *memcache.Client) *SessionManager {
	return &SessionManager{client: client}
}

func (m *SessionManager) CreateSession(user user.User) (string, error) {
	id := uuid.NewV4()
	tokenString, err := createJWT(id.String(), user)
	if err != nil {
		return "", fmt.Errorf("create token failed: %s", err)
	}
	session := Session{ID: id.String(), User: user}

	item := memcache.Item{Key: session.ID, Value: []byte(strconv.Itoa(int(user.ID))), Expiration: ExpireTimeDay * 24 * 3600}

	err = m.client.Set(&item)
	if err != nil {
		return "", fmt.Errorf("save session failed: %s", err)
	}
	return tokenString, nil
}

func (m *SessionManager) GetSession(token string) (*Session, error) {

	token_, err := jwt.ParseWithClaims(token, &Claims{}, hashSecretGetter)
	if err != nil {
		return nil, fmt.Errorf("token parser failed")
	}
	claims, ok := token_.Claims.(*Claims)

	if !ok || !token_.Valid {
		return nil, fmt.Errorf("token is invalid")
	}
	item, err := m.client.Get(claims.SessionID)
	if err != nil {
		return nil, fmt.Errorf("get session with id: %s failed: %s", claims.SessionID, err)
	}
	equal := reflect.DeepEqual(item.Value, []byte(strconv.Itoa(int(claims.User.ID))))
	if !equal {
		return nil, fmt.Errorf("get session with id: %s failed: user id does not match", claims.SessionID)
	}

	session := &Session{claims.SessionID, claims.User}
	return session, nil
}

//TODO должен вызываться при LOGOUT
func (m *SessionManager) DestroySession(sessionId string) error {
	if err := m.client.Delete(sessionId); err != nil {
		return fmt.Errorf("destroy session failed: %s", err)
	}
	return nil
}

func hashSecretGetter(token *jwt.Token) (interface{}, error) {
	method, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok || method.Alg() != "HS256" {
		return nil, fmt.Errorf("bad sign method")
	}
	//TODO!!!!!
	//secretKey, ok := os.LookupEnv("SECRET_KEY")
	//if !ok {
	//	return "", fmt.Errorf("SECRET_KEY is required")
	//}
	return []byte(SecretKey), nil
}
