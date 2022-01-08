package middlewares

import (
	"net/http"
	"strings"

	"github.com/solabsafrica/afrikanest/exceptions"
	"github.com/solabsafrica/afrikanest/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AuthChecker interface {
	Check(*gin.Context)
}

type authCheckImpl struct {
	authService service.AuthServiceWithContext
}

func (checker *authCheckImpl) Check(ctx *gin.Context) {
	context := ctx.Request.Context()
	authorization := ctx.GetHeader("Authorization")
	parts := strings.Split(authorization, " ")
	if len(parts) < 2 {
		ctx.JSON(http.StatusUnauthorized, exceptions.AuthFailed)
		ctx.Abort()
		return
	}
	token, err := checker.authService(context).Validate(parts[1], jwt.MapClaims{})
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err)
		ctx.Abort()
		return
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	ctx.Set("identity", claims["sub"])
}

func NewAuthChecker(authService service.AuthServiceWithContext) AuthChecker {
	return &authCheckImpl{
		authService: authService,
	}
}
