package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Generate tokens
// @Tags Auth
// @Description Generate tokens
// @ID generate-tokens
// @Produce  json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} transort_error
// @Failure 500 {object} transort_error
// @Failure default {object} transort_error
// @Router /auth/getTokens/ [get]
func (h *Handler) getTokens(c *gin.Context) {
	guid := c.Param("guid")

	ip := c.GetHeader("Ip")

	accessToken, refreshToken, err := h.services.GenerateTokens(guid, ip)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

// @Summary Refresh
// @Tags Auth
// @Description Refresh
// @ID refresh
// @Produce  json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} transort_error
// @Failure 500 {object} transort_error
// @Failure default {object} transort_error
// @Router /auth/refresh/ [post]
func (h *Handler) refresh(c *gin.Context) {

}

// @Summary Register
// @Tags Auth
// @Description Register
// @ID register
// @Produce  json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} transort_error
// @Failure 500 {object} transort_error
// @Failure default {object} transort_error
// @Router /auth/register/ [post]
func (h *Handler) register(c *gin.Context) {

}
