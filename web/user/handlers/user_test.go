package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/artkescha/checker/online_checker/pkg/session"
	"github.com/artkescha/checker/online_checker/pkg/user"
	"github.com/artkescha/checker/online_checker/pkg/user/repository"
	"github.com/artkescha/checker/online_checker/web/request"
	"github.com/golang/mock/gomock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type RequestBuilder struct{ host string }

type TestClient struct {
	// токен, по которому происходит авторизация на внешней системе, уходит туда через хедер
	accessToken string
	// урл внешней системы, куда идти
	url     string
	handler func(http.ResponseWriter, *http.Request)
}

type ErrorResponse struct {
	Error string
}

func (r RequestBuilder) RegisterUser(params request.Login) (*http.Request, error) {
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", r.host+"/api/register", bytes.NewBuffer(jsonParams))
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (r RequestBuilder) LoginUser(params request.Login) (*http.Request, error) {
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", r.host+"/api/login", bytes.NewBuffer(jsonParams))
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (c *TestClient) TestUser(request *http.Request) (string, error) {
	resp := httptest.NewRecorder()
	c.handler(resp, request)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read all response err: %s", err)
	}

	if resp.Code != http.StatusOK && resp.Code != http.StatusCreated {
		errResp := ErrorResponse{}
		err = json.Unmarshal(body, &errResp)
		if err != nil {
			return "", fmt.Errorf("cant unpack result json: %s", err)
		}
		return "", fmt.Errorf(errResp.Error)
	}
	var token = struct {
		Value string `json:"token"`
	}{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return "", fmt.Errorf("cant unpack result json: %s", err)
	}
	return token.Value, nil
}

func TestUserHandler_RegisterUser(t *testing.T) {

	tests := []struct {
		name    string
		params  request.Login
		user    *user.User
		want    string
		wantErr bool
	}{
		{
			name: "ok_test empty data",
			params: request.Login{
				Username: "",
				Password: "",
			},
			user:    &user.User{Name: "", Password: ""},
			want:    "",
			wantErr: false,
		},
		{
			name: "ok_test",
			params: request.Login{
				Username: "1",
				Password: "test",
			},
			user:    &user.User{Name: "1", Password: "test"},
			want:    "1",
			wantErr: false,
		},
		{
			name: "internal server status",
			params: request.Login{
				Username: "1",
				Password: "test",
			},
			user:    &user.User{},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		defer server.Close()

		request, err := RequestBuilder{server.URL}.RegisterUser(tt.params)
		if err != nil {
			continue
		}

		ctrl := gomock.NewController(t)
		mockRepo := repository.NewMockUserRepo(ctrl)
		mockManager := session.NewMockManager(ctrl)

		if tt.wantErr {
			err := fmt.Errorf("internal server error")
			mockRepo.EXPECT().Insert(tt.params.Username, user.GetMD5Password(tt.params.Password)).Return(tt.user, err).AnyTimes()
			mockManager.EXPECT().CreateSession(*tt.user).Return("", err).AnyTimes()
		} else if tt.params.Username == "" && tt.params.Password == "" {
			mockRepo.EXPECT().Insert(tt.params.Username, user.GetMD5Password(tt.params.Password)).Return(tt.user, err).AnyTimes()
			mockManager.EXPECT().CreateSession(*tt.user).Return("", err).AnyTimes()
		} else {
			mockRepo.EXPECT().Insert(tt.params.Username, user.GetMD5Password(tt.params.Password)).Return(tt.user, nil).AnyTimes()
			mockManager.EXPECT().CreateSession(*tt.user).Return("1", nil).AnyTimes()
		}

		userHandlers := UserHandler{
			UsersRepo:      mockRepo,
			SessionManager: mockManager,
		}
		client := TestClient{
			accessToken: "",
			url:         server.URL,
			handler:     userHandlers.Register,
		}

		got, err := client.TestUser(request)
		if (err != nil) != tt.wantErr {
			t.Errorf("RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
			return
		}

		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got: %+v tt.want: %+v", got, tt.want)
			return
		}
	}
}

func TestUserHandler_LoginUser(t *testing.T) {
	tests := []struct {
		name     string
		params   request.Login
		user     *user.User
		want     string
		wantErr  bool
	}{
		{
			name: "ok_test empty data",
			params: request.Login{
				Username: "",
				Password: "",
			},
			user:    &user.User{Name: "", Password: ""},
			want:    "",
			wantErr: false,
		},
		{
			name: "ok_test",
			params: request.Login{
				Username: "1",
				Password: "test",
			},
			user:    &user.User{Name: "1", Password: user.GetMD5Password("test")},
			want:    "1",
			wantErr: false,
		},
		{
			name: "internal server status",
			params: request.Login{
				Username: "1",
				Password: "test",
			},
			user:    &user.User{},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		defer server.Close()

		request, err := RequestBuilder{server.URL}.LoginUser(tt.params)
		if err != nil {
			continue
		}

		ctrl := gomock.NewController(t)
		mockRepo := repository.NewMockUserRepo(ctrl)
		mockManager := session.NewMockManager(ctrl)

		if tt.wantErr {
			err := fmt.Errorf("internal server error")
			mockRepo.EXPECT().GetUserByLogin(tt.params.Username).Return(tt.user, err).AnyTimes()
			mockManager.EXPECT().CreateSession(*tt.user).Return("", err).AnyTimes()
		} else if tt.params.Username == "" && tt.params.Password == "" {
			mockRepo.EXPECT().GetUserByLogin(tt.params.Username).Return(tt.user, err).AnyTimes()
			mockManager.EXPECT().CreateSession(*tt.user).Return("", err).AnyTimes()
		} else {
			mockRepo.EXPECT().GetUserByLogin(tt.params.Username).Return(tt.user, nil).AnyTimes()
			mockManager.EXPECT().CreateSession(*tt.user).Return("1", nil).AnyTimes()
		}

		userHandlers := UserHandler{
			UsersRepo:      mockRepo,
			SessionManager: mockManager,
		}
		client := TestClient{
			accessToken: "",
			url:         server.URL,
			handler:     userHandlers.Login,
		}

		got, err := client.TestUser(request)
		if (err != nil) != tt.wantErr {
			t.Errorf("LoginUser() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got: %+v tt.want: %+v", got, tt.want)
			return
		}
	}
}
