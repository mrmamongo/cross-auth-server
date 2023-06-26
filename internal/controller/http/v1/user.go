package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mrmamongo/go-auth-tg/internal/entity"
	"github.com/mrmamongo/go-auth-tg/internal/usecase"
	"github.com/mrmamongo/go-auth-tg/pkg/logger"
	"net/http"
)

type userRoutes struct {
	t usecase.User
	l logger.Interface
}

func newUserRoutes(handler *gin.RouterGroup, t usecase.User, l logger.Interface) {
	r := &userRoutes{t, l}

	h := handler.Group("/user")
	{
		h.GET("/", r.getAll)
		h.GET("/:username", r.getByUsername)
		h.POST("/", r.create)
		h.PUT("/:username", r.update)
	}
}

// @Summary     Users
// @Description Get all users
// @ID          user
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Success     200 {object} []entity.User
// @Failure     500 {object} response
// @Router      /user [get]
func (r *userRoutes) getAll(c *gin.Context) {
	users, err := r.t.GetAll(c.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - getAll")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// @Summary     User
// @Description Get user by username
// @ID          user_by_username
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Param       username path string true "Username"
// @Success     200 {object} []entity.User
// @Failure     404 {object} response
// @Router      /user/:username [get]
func (r *userRoutes) getByUsername(c *gin.Context) {
	username := c.Param("username")
	user, err := r.t.GetByUsername(c.Request.Context(), username)
	if err != nil {
		r.l.Error(err, "http - v1 - getByUsername")
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

type userCreateForm struct {
	Username         string `json:"username" binding:"required"`
	TelegramUsername string `json:"telegram_username" binding:"required"`
}

// @Summary     New User
// @Description Create new user
// @ID          user_create
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Param       request body userCreateForm true "User Create Form"
// @Success     200 {object} entity.User
// @Failure     409 {object} response
// @Failure     500 {object} response
// @Router      /user [post]
func (r *userRoutes) create(c *gin.Context) {
	var user userCreateForm
	if err := c.ShouldBindJSON(&user); err != nil {
		r.l.Error(err, "http - v1 - create")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created := entity.User{
		Username:         user.Username,
		TelegramUsername: user.TelegramUsername,
	}
	if err := r.t.Create(c.Request.Context(), &created); err != nil {
		r.l.Error(err, "http - v1 - create")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

type userUpdateRequest struct {
	Username         string `json:"username" binding:"omitempty"`
	TelegramUsername string `json:"telegramUsername" binding:"omitempty"`
}

// @Summary     Update User
// @Description Update user
// @ID          user_update
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Param       request query userUpdateRequest true "User Update Request"
// @Success     200 {object} entity.User
// @Failure     409 {object} response
// @Failure     500 {object} response
// @Router      /user/:username [post]
func (r *userRoutes) update(c *gin.Context) {
	username := c.Param("username")
	var req userUpdateRequest
	var err error
	if err = c.ShouldBindQuery(&req); err != nil {
		r.l.Error(err, "http - v1 - update")
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var user *entity.User
	if user, err = r.t.GetByUsername(c.Request.Context(), username); err != nil {
		r.l.Error(err, "http - v1 - update")
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	user.Username = req.Username
	user.TelegramUsername = req.TelegramUsername

	if err := r.t.Update(c.Request.Context(), user); err != nil {
		r.l.Error(err, "http - v1 - update")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
