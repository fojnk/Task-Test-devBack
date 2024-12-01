package transport

import (
	"net/http"

	"github.com/fojnk/Task-Test-devBack/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TokenPair struct {
	AccesToken   string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type InputRegister struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary Generate tokens
// @Tags Auth
// @Description Generate tokens
// @ID generate-tokens
// @Param guid query string true "Guid"
// @Param Ip header string true "Ip"
// @Produce  json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} transort_error
// @Failure 500 {object} transort_error
// @Failure default {object} transort_error
// @Router /auth/getTokens [get]
func (h *Handler) getTokens(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	guids, ok := queryParams["guid"]
	if !ok {
		NewTransportErrorResponse(c, http.StatusBadRequest, "bad param format")
		return
	}
	guid := guids[0]

	logrus.Info(guid)

	ip := c.GetHeader("Ip")

	accessToken, refreshToken, err := h.services.GenerateTokens(guid, ip)
	if err != nil {
		NewTransportErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

// @Summary Refresh
// @Param Ip header string true "Ip"
// @Tags Auth
// @Description Refresh
// @ID refresh
// @Accept  json
// @Produce  json
// @Param input body TokenPair true "tokens"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} transort_error
// @Failure 500 {object} transort_error
// @Failure default {object} transort_error
// @Router /auth/refresh [post]
func (h *Handler) refresh(c *gin.Context) {
	header := c.GetHeader("Ip")

	var input TokenPair

	if err := c.BindJSON(&input); err != nil {
		NewTransportErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	newAccess, newRefresh, err := h.services.Refresh(input.AccesToken, input.RefreshToken, header)
	if err != nil {
		NewTransportErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken":  newAccess,
		"refreshToken": newRefresh,
	})
}

// @Summary Register
// @Tags Auth
// @Description create account
// @ID create-account
// @Param Ip header string true "Ip"
// @Accept  json
// @Produce  json
// @Param input body InputRegister true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404,500,default {object} transort_error
// @Router /auth/register [post]
func (h *Handler) register(c *gin.Context) {
	header := c.GetHeader("Ip")
	var input InputRegister

	if err := c.BindJSON(&input); err != nil {
		NewTransportErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user := models.User{
		Email:    input.Email,
		Password: input.Password,
		Username: input.Username,
	}

	logrus.Printf("create user with %s, %s", input.Email, input.Username)

	id, err := h.services.Authorization.CreateUser(user)
	if err != nil {
		NewTransportErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Printf("generate tokens for user: %s", id)

	accessToken, refreshToken, err := h.services.Authorization.GenerateTokens(id, header)
	if err != nil {
		NewTransportErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"Guid":         id,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}
