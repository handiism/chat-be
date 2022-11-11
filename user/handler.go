package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/handiism/chat/utils"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{
		service: service,
	}
}

func (h *Handler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := LoginRequest{}
		err := c.ShouldBindJSON(&request)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.NewFailResponse(err.Error(), nil))
			return
		}

		user, err := h.service.Login(request.Username, request.Password)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, utils.NewFailResponse(err.Error(), nil))
			return
		}

		c.JSON(http.StatusOK, utils.NewSuccessResponse("Login success.", user))
	}
}

func (h *Handler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := RegisterRequest{}
		err := c.BindJSON(&request)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.NewFailResponse(err.Error(), nil))
			return
		}

		user, err := h.service.Register(request.Username, request.Password)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, utils.NewFailResponse(err.Error(), nil))
			return
		}

		c.JSON(http.StatusCreated, utils.NewSuccessResponse("User registered", user))
	}
}

func (h *Handler) Fetch() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if len(id) == 0 {
			c.JSON(http.StatusUnprocessableEntity, utils.NewFailResponse("User ID not found", nil))
			return
		}

		user, err := h.service.FetchById(id)
		if err != nil {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
