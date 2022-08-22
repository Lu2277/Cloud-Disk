package middleware

import (
	"Cloud-Disk/core/helper"
	"net/http"
)

type AuthMiddleware struct {
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
		userClaims, err := helper.AnalyseToken(auth)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
		//在Header头部设置认证的相关信息
		r.Header.Set("UserId", string(rune(userClaims.ID)))
		r.Header.Set("UserIdentity", userClaims.Identity)
		r.Header.Set("UserName", userClaims.Name)
		next(w, r)
	}
}
