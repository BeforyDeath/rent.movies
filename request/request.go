package request

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

// todo нигде не используется, хмм
func GetRouterParams(r *http.Request) httprouter.Params {
	ctx := r.Context()
	return ctx.Value("params").(httprouter.Params)
}

func GetClaims(r *http.Request) jwt.MapClaims {
	ctx := r.Context()
	return ctx.Value("claims").(jwt.MapClaims)
}

func GetJson(r *http.Request) map[string]interface{} {
	ctx := r.Context()
	return ctx.Value("request").(map[string]interface{})
}
