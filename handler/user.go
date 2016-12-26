package handler

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/BeforyDeath/rent.movies/request"
	"github.com/BeforyDeath/rent.movies/storage"
	"github.com/BeforyDeath/rent.movies/validator"
	jwt "github.com/dgrijalva/jwt-go"
)

type user struct{}

func (u user) Login(w http.ResponseWriter, r *http.Request) {
	res := Result{}
	defer func() { json.NewEncoder(w).Encode(res) }()

	req := request.GetJson(r)

	user := new(storage.User)
	err := validator.GetRequest(req, user)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user.Pass, _ = u.Hashed(user.Pass, cfg.Security.PasswordSalt)
	err = user.Check()
	if err != nil {
		res.Error = err.Error()
		return
	}

	token, err := u.CreateToken([]byte(cfg.Security.TokenSalt), user)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Success = true
	res.Data = struct{ Token string }{Token: token}

	return
}

func (u user) Create(w http.ResponseWriter, r *http.Request) {
	res := Result{}
	defer func() { json.NewEncoder(w).Encode(res) }()

	req := request.GetJson(r)

	user := new(storage.User)
	err := validator.GetRequest(req, user)
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user.Pass, _ = u.Hashed(user.Pass, cfg.Security.PasswordSalt)
	user.CreateAt = time.Now()

	err = user.Create()
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//res.Data = user
	res.Success = true
	w.WriteHeader(http.StatusCreated)
	return
}

func (u user) Authorization(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		req := request.GetJson(r)

		if token, ok := req["token"]; ok && reflect.ValueOf(token).Kind() == reflect.String && token != "" {

			claims, err := u.CheckToken(token.(string), cfg.Security.TokenSalt)
			if err != nil {
				res := Result{Error: err.Error()}
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(res)
				return
			}

			ctx := context.WithValue(r.Context(), "claims", claims)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
			return

		}

		res := Result{Error: "One of the parameters specified was missing or invalid: address is taken"}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(res)

		return
	}
	return http.HandlerFunc(fn)
}

func (u user) Hashed(in, salt string) (string, error) {
	mac := hmac.New(sha256.New, []byte(salt))
	_, err := mac.Write([]byte(in))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", mac.Sum(nil)), nil
}

// todo нуно подмешать заголовки и устройства входа, от угона токена
func (u user) CreateToken(mySigningKey []byte, su *storage.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := make(jwt.MapClaims)
	claims["user_id"] = su.Id
	claims["login"] = su.Login
	//claims["Name"] = u.Name
	//claims["Age"] = u.Age
	//claims["Phone"] = u.Phone
	//claims["CreateAt"] = u.CreateAt
	claims["exp"] = time.Now().Add(time.Second * cfg.Security.TokenExpired).Unix()
	token.Claims = claims

	tokenString, err := token.SignedString(mySigningKey)
	return tokenString, err
}

func (u user) CheckToken(myToken string, myKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(myKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("This token is terrible!")
}
