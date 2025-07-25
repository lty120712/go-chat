package controller

import (
	"github.com/gin-gonic/gin"
	interfacesservice "go-chat/internal/interfaces/service"
	request "go-chat/internal/model/request"
	"strconv"
)

// GroupController 群组相关控制器
// @Tags Group
// @Description 群组相关控制器
type GroupController struct {
	BaseController
	groupService interfacesservice.GroupServiceInterface
}

var GroupControllerInstance *GroupController

func InitGroupController(groupService interfacesservice.GroupServiceInterface) {
	GroupControllerInstance = &GroupController{
		groupService: groupService,
	}
}

// Create 创建群组
// @Summary 创建群组
// @Description 创建一个新的群组，并将指定的成员添加到群组中
// @Tags Group
// @Accept json
// @Produce json
// @Param group body model.GroupCreateRequest true "创建群组请求参数"
// @Success 200 {object} model.Response "群组创建成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 500 {object} model.Response "内部服务器错误"
// @Router /group/create [post]
func (con GroupController) Create(c *gin.Context) {
	var req *request.GroupCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}
	err := req.Validate()
	if err != nil {
		con.Error(c, err.Error())
		return
	}

	if err := con.groupService.Create(req); err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c)
}

func (con GroupController) Update(c *gin.Context) {

}

func (con GroupController) Join(c *gin.Context) {
	groupIdStr, _ := c.GetQuery("group_id")
	groupId, _ := strconv.ParseUint(groupIdStr, 10, 64)
	userId, ok := c.Get("id")
	if !ok {
		con.Error(c, "need user_id")
		return
	}

	if err := con.groupService.Join(uint(groupId), userId.(uint)); err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c)
}

func (con GroupController) Quit(c *gin.Context) {
	groupIdStr, _ := c.GetQuery("group_id")
	groupId, _ := strconv.ParseUint(groupIdStr, 10, 64)
	userId, ok := c.Get("id")
	if !ok {
		con.Error(c, "need user_id")
		return
	}
	if err := con.groupService.Quit(uint(groupId), userId.(uint)); err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c)
}

func (con GroupController) Search(c *gin.Context) {

}

func (con GroupController) Member(c *gin.Context) {
	groupIdStr := c.Query("group_id")
	groupId, _ := strconv.ParseUint(groupIdStr, 10, 64)
	members, err := con.groupService.Member(uint(groupId))
	if err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c, members)
}
