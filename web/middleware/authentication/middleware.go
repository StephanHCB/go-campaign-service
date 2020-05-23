package authentication

import (
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

// TODO we need asymmetric crypto because we want to keep the key only in the IAM,
//      also frontends should be able to validate tokens

func createAndConfigureAuthenticationMiddleware(secret string) *jwtmiddleware.JWTMiddleware {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
		// Allow missing credentials, will leave the "user" context key unset (which you should interpret as "not authenticated")
		CredentialsOptional: true,
	})
	return jwtMiddleware
}

func createHandlerFunction(authMw *jwtmiddleware.JWTMiddleware, next http.Handler) func(w http.ResponseWriter, r *http.Request) {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		err := authMw.CheckJWT(w, r)
		if err != nil {
			// TODO react to error with some more detail than just a 401
			// note that this error does not trigger if the Authorization header is missing completely, only if
			// there is something wrong with it
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
	return handlerFunc
}

func AddJWTTokenInfoToContext(secret string) func(http.Handler) http.Handler {
	authMw := createAndConfigureAuthenticationMiddleware(secret)

	middlewareCreator := func(next http.Handler) http.Handler {
		fn := createHandlerFunction(authMw, next)
		return http.HandlerFunc(fn)
	}

	return middlewareCreator
}
