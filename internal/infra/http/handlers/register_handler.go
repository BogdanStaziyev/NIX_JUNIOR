package handlers

import (
	"fmt"
	"github.com/BogdanStaziyev/NIX_Junior/internal/app"
	"github.com/BogdanStaziyev/NIX_Junior/internal/domain"
	"github.com/BogdanStaziyev/NIX_Junior/internal/infra/http/requests"
	"github.com/BogdanStaziyev/NIX_Junior/internal/infra/http/response"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type RegisterHandler struct {
	as app.AuthService
}

func NewRegisterHandler(a app.AuthService) RegisterHandler {
	return RegisterHandler{
		as: a,
	}
}

// Register 		godoc
// @Summary 		Register
// @Description 	New user registration
// @ID				user-register
// @Tags			Auth Actions
// @Accept 			json
// @Produce 		json
// @Param			input body requests.RegisterAuth true "users email, users password"
// @Success 		201 {object} response.UserResponse
// @Failure			400 {object} response.Error
// @Router			/register [post]
func (r RegisterHandler) Register(ctx echo.Context) error {
	var registerUser requests.RegisterAuth
	if err := ctx.Bind(&registerUser); err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Could not decode user data")
	}
	if err := ctx.Validate(&registerUser); err != nil {
		return response.ErrorResponse(ctx, http.StatusUnprocessableEntity, "Could not validate user data")
	}

	userFromRegister := registerUser.RegisterToUser()

	user, err := r.as.Register(userFromRegister)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Could not save new user: %s", err))
	}
	userResponse := domain.User.DomainToResponse(user)
	return response.Response(ctx, http.StatusCreated, userResponse)
}

// Login 			godoc
// @Summary 		LoginAuth
// @Description 	LoginAuth
// @Tags			Auth Actions
// @Accept 			json
// @Produce 		json
// @Param			input body requests.LoginAuth true "users email, users password"
// @Success 		201 {object} response.LoginResponse
// @Failure			400 {object} response.Error
// @Failure			401 {object} response.Error
// @Failure			500 {object} response.Error
// @Router			/login [post]
func (r RegisterHandler) Login(ctx echo.Context) error {
	var authUser requests.LoginAuth
	if err := ctx.Bind(&authUser); err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Could not decode user data")
	}
	if err := ctx.Validate(&authUser); err != nil {
		return response.ErrorResponse(ctx, http.StatusUnprocessableEntity, "Could not validate user data")
	}
	accessToken, refreshToken, exp, err := r.as.Login(authUser)
	if err != nil {
		if strings.HasSuffix(err.Error(), "upper: no more rows in this result set") {
			return response.ErrorResponse(ctx, http.StatusNotFound, fmt.Sprintf("Could not login, user not exists: %s", err))
		} else {
			return response.ErrorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Could not login user: %s", err))
		}
	}
	res := response.NewLoginResponse(accessToken, refreshToken, exp)
	return response.Response(ctx, http.StatusOK, res)
}
