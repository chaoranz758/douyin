package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
)

func CSRFMiddle() gin.HandlerFunc {
	csrfMiddleware := csrf.Protect([]byte("32-byte-long-auth-key"))
	return adapter.Wrap(csrfMiddleware)
}
