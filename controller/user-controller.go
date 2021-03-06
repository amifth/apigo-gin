package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/amifth/apigo-gin/dto"
	"github.com/amifth/apigo-gin/helper"
	"github.com/amifth/apigo-gin/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	Update(context *gin.Context)
	Profile(context *gin.Context)
	AllUser(context *gin.Context)
}

type userController struct {
	userService service.UserService
	jwtService  service.JWTService
}

func NewUserController(userService service.UserService, jwtService service.JWTService) UserController {
	return &userController{
		userService: userService,
		jwtService:  jwtService,
	}
}

//GetUsers ... Get all users
// func (c *userController) GetUsers(context *gin.Context) {
// 	var user dto.UserUpdateDTO
// 	err := Models.GetAllUsers(&user)
// 	if err != nil {
// 		context.AbortWithStatus(http.StatusNotFound)
// 	} else {
// 		context.JSON(http.StatusOK, user)
// 	}
// }

func (c *userController) Update(context *gin.Context) {
	var userUpdateDTO dto.UserUpdateDTO
	errDTO := context.ShouldBind(&userUpdateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)

	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	id, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		panic(err.Error())
	}
	userUpdateDTO.ID = id
	u := c.userService.Update(userUpdateDTO)
	res := helper.BuildResponse("200", true, "Successful!", u)
	context.JSON(http.StatusOK, res)
}

// User godoc
// @Summary      user account
// @Description  user update
// @Tags         user
// @Accept       json
// @Produce      json
// @Param 		 Authorization header string true "Bearer"
// @Param        name    query     string  false  "name"  Format(name)
// @Param        email    query     string  false  "email"  Format(email)
// @Param        password    query     string  true  "password"  Format(password)
// @Success      200  {object}  map[string]interface{}
// @Router       /user/profile [put]
func (c *userController) Profile(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	token, err := c.jwtService.ValidateToken(authHeader)
	if err != nil {
		panic(err.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	user := c.userService.Profile(id)
	res := helper.BuildResponse("200", true, "Successful!", user)
	context.JSON(http.StatusOK, res)
}

// User godoc
// @Summary      user all account
// @Description  user get all account
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /user/all [get]
func (c *userController) AllUser(context *gin.Context) {
	users := c.userService.AllUser()
	response := helper.BuildResponse("200", true, "Successful!", users)
	context.JSON(http.StatusOK, response)
}
