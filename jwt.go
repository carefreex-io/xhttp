package xhttp

import (
	"errors"
	"github.com/carefreex-io/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"sync"
)

var (
	ErrEmptyToken   = errors.New("not found token in header")
	ErrTokenInvalid = errors.New("token is invalid")
)

type XJwt struct {
	Secret    []byte
	LookupKey string
	Auth      Authorization
}

type JwtOptions struct {
	Secret string
	Lookup string
	Auth   Authorization
}

var (
	xJwt     *XJwt
	xJwtOnce sync.Once
)

type Authorization interface {
	Verify(claims jwt.Claims) bool
	GetEmptyMyClaims() jwt.Claims
}

func NewXJwt(options JwtOptions) *XJwt {
	if xJwt == nil {
		xJwtOnce.Do(func() {
			xJwt = &XJwt{
				Secret:    []byte(options.Secret),
				LookupKey: options.Lookup,
				Auth:      options.Auth,
			}
		})
	}

	return xJwt
}

func (j *XJwt) GenerateToken(claims jwt.Claims) (sign string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(j.Secret)
}

func (j *XJwt) Middleware(exception map[string]byte) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _, ok := exception[ctx.Request.URL.Path]; ok {
			ctx.Next()
			return
		}

		tokenStr := ctx.Request.Header.Get(xJwt.LookupKey)
		if tokenStr == "" {
			UnauthorizedResponse(ctx, ErrEmptyToken)
			ctx.Abort()
			return
		}

		claims := xJwt.Auth.GetEmptyMyClaims()

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return xJwt.Secret, nil
		})
		if err != nil {
			logger.ErrorfX(ctx, "jwt.ParseWithClaims failed: tokenStr=%v claims=%v err=%v", tokenStr, claims, err)
			UnauthorizedResponse(ctx, err)
			ctx.Abort()
			return
		}

		if token.Valid && !xJwt.Auth.Verify(token.Claims) {
			logger.ErrorfX(ctx, "token verify failed: token.Valid=%v token.Claims=%v", token.Valid, token.Claims)
			UnauthorizedResponse(ctx, ErrTokenInvalid)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
