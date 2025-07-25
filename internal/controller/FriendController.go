package controller

import (
	"github.com/gin-gonic/gin"
	interfacesservice "go-chat/internal/interfaces/service"
	request "go-chat/internal/model/request"
	"strconv"
)

// FriendController 用户相关控制器
// @Tags Friend
// @Description 控制好友相关的 API
type FriendController struct {
	BaseController
	friendService interfacesservice.FriendServiceInterface
}

var FriendControllerInstance *FriendController

func InitFriendController(friendService interfacesservice.FriendServiceInterface) {
	FriendControllerInstance = &FriendController{
		friendService: friendService,
	}
}

func (con *FriendController) Add(c *gin.Context) {
	id := c.GetUint("id")
	var req []uint
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}
	if err := con.friendService.Add(id, req); err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c)
}
func (con *FriendController) ListReq(c *gin.Context) {
	id := c.GetUint("id")
	data, err := con.friendService.ListReq(id)
	if err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c, data)
}
func (con *FriendController) HandleReq(c *gin.Context) {

	var req request.FriendHandlerReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		con.Error(c, err.Error())
		return
	}
	if err := con.friendService.HandleReq(req.Id, req.Status); err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c)
}
func (con *FriendController) Remove(c *gin.Context) {
	id := c.GetUint("id")
	var req []int64
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}
	if err := con.friendService.Remove(id, req); err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c)
}

func (con *FriendController) GroupCreate(c *gin.Context) {
	id := c.GetUint("id")
	var req request.FriendGroupCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}
	if err := con.friendService.GroupCreate(id, req); err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c)
}

func (con *FriendController) GroupDelete(c *gin.Context) {
	groupIdStr := c.Query("id")
	groupId, _ := strconv.ParseInt(groupIdStr, 10, 64)
	if err := con.friendService.GroupDelete(groupId); err != nil {
		con.Error(c, err.Error())
		return
	}
}

func (con *FriendController) GroupUpdate(c *gin.Context) {
	var req request.FriendGroupUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}

	if err := con.friendService.GroupUpdate(req); err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c)
}
func (con *FriendController) GroupList(c *gin.Context) {
	userId := c.GetUint("id")

	data, err := con.friendService.GroupList(userId)
	if err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c, data)
}
