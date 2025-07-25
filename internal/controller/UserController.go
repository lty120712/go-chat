package controller

import (
	"github.com/gin-gonic/gin"
	interfacesservice "go-chat/internal/interfaces/service"
	"go-chat/internal/model"
	request "go-chat/internal/model/request"
	"go-chat/internal/service"
	"strconv"
)

// UserController 用户相关控制器
// @Tags User
// @Description 控制用户相关的 API，包含用户注册、登录和 ping 等接口
type UserController struct {
	BaseController
	userService interfacesservice.UserServiceInterface
}

var UserControllerInstance *UserController

func InitUserController(userService interfacesservice.UserServiceInterface) {
	UserControllerInstance = &UserController{
		userService: userService,
	}
}

// Ping 返回一个 ping 的响应，测试接口是否可用
// @Summary Ping 接口
// @Description 返回一个简单的 ping 响应，用于检测 API 是否可用
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} model.Response "成功"
// @Router /user/ping [get]
func (con UserController) Ping(c *gin.Context) {
	// 调用服务层获取数据
	//返回结果
	con.Success(c, "pong ping ping pang")
}

// Register 用户注册接口
// @Summary 用户注册
// @Description 接口用于用户注册，提供注册所需的字段
// @Tags user
// @Accept json
// @Produce json
// @Param username body string true "用户名"
// @Param password body string true "密码"
// @Success 200 {object} model.Response "成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Router /user/register [post]
func (con UserController) Register(c *gin.Context) {
	registerRequest := &request.RegisterRequest{}
	if err := c.ShouldBindJSON(registerRequest); err != nil {
		con.Error(c, err.Error())
		return
	}
	if err := service.UserServiceInstance.Register(registerRequest.Username, registerRequest.Password, registerRequest.RePassword); err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c)
	return
}

// Login 用户登录接口
// @Summary 用户登录
// @Description 用户登录，返回用户的认证信息
// @Tags user
// @Accept json
// @Produce json
// @Param body body model.LoginRequest true "登录信息"
// @Success 200 {object} model.Response "成功"
// @Failure 401 {object} model.Response "登陆失败"
// @Router /user/login [post]
func (con UserController) Login(c *gin.Context) {
	loginRequest := &request.LoginRequest{}
	if err := c.ShouldBindJSON(loginRequest); err != nil {
		con.Error(c, err.Error())
		return
	}
	token, err := con.userService.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c, token)
	return
}

func (con UserController) Logout(c *gin.Context) {
	id := c.GetUint("id")
	con.userService.Logout(id)
	con.Success(c)
	return
}

func (con UserController) OnlineStatusChange(c *gin.Context) {
	id := c.GetUint("id")
	onlineStatus, _ := strconv.ParseInt(c.Query("online_status"), 10, 64)
	if err := con.userService.OnlineStatusChange(id, model.OnlineStatus(onlineStatus)); err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c)
}

func (con UserController) GetUserInfo(c *gin.Context) {
	id := c.GetUint("id")
	userInfo, err := con.userService.GetUserInfo(id)
	if err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c, userInfo)
}

// Update 用户更新信息接口
// @Summary 更新用户信息
// @Description 需要登录
// @security Bearer
// @Tags user
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "用户信息"
// @Success 200 {string} string "成功"
// @Failure 401 {string} string "未授权"
// @Router /user/update [post]
func (con UserController) Update(c *gin.Context) {
	updateRequest := &request.UserUpdateRequest{}
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		con.Error(c, err.Error())
		return
	}
	err := updateRequest.Validate()
	if err != nil {
		con.Error(c, err.Error())
		return
	}
	if err := con.userService.UpdateUser(updateRequest); err != nil {
		con.Error(c, err)
		return
	}
	con.Success(c)
}
