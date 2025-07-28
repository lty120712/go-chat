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

// Join 加入群组
// @Summary 加入群组
// @Description 加入群组
// @Tags Group
// @Accept json
// @Produce json
// @Param group_id query uint true "群组ID"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /group/join [Get]
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

// Quit 退出群组
// @Summary 退出群组
// @Description 退出群组
// @Tags Group
// @Accept json
// @Produce json
// @Param group_id query uint true "群组ID"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /group/quit [Get]
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
	var req request.GroupSearchRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		con.Error(c, err.Error())
		return
	}

	resp, err := con.groupService.Page(req)
	if err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c, resp)
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

// CreateAnnouncement 创建群组公告
// @Summary 创建群组公告
// @Description 创建一个新的群组公告
// @Tags Group
// @Accept json
// @Produce json
// @Param group_id path int true "群组ID"
// @Param data body request.GroupAnnouncementCreateRequest true "公告内容"
// @Success 200 {object} model.Response
// @Router /group/{group_id}/announcement/create [post]
func (con GroupController) CreateAnnouncement(c *gin.Context) {
	groupIdStr := c.Param("group_id")
	groupId, _ := strconv.Atoi(groupIdStr)

	var req request.GroupAnnouncementCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}

	if err := con.groupService.CreateAnnouncement(uint(groupId), &req); err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c)
}

// UpdateAnnouncement 更新群组公告
// @Summary 更新群组公告
// @Description 更新指定的群组公告
// @Tags Group
// @Accept json
// @Produce json
// @Param group_id path int true "群组ID"
// @Param data body request.GroupAnnouncementUpdateRequest true "公告更新内容"
// @Success 200 {object} model.Response
// @Router /group/{group_id}/announcement/update [post]
func (con GroupController) UpdateAnnouncement(c *gin.Context) {
	groupIdStr := c.Param("group_id")
	groupId, _ := strconv.Atoi(groupIdStr)

	var req request.GroupAnnouncementUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}

	if err := con.groupService.UpdateAnnouncement(uint(groupId), &req); err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c)
}

// DeleteAnnouncement 删除群组公告
// @Summary 删除群组公告
// @Description 删除指定群组的公告
// @Tags Group
// @Accept json
// @Produce json
// @Param group_id path int true "群组ID"
// @Param announcement_id query int true "公告ID"
// @Success 200 {object} model.Response
// @Router /group/{group_id}/announcement/delete [delete]
func (con GroupController) DeleteAnnouncement(c *gin.Context) {
	groupIdStr := c.Param("group_id")
	groupId, _ := strconv.Atoi(groupIdStr)

	announcementIdStr := c.DefaultQuery("announcement_id", "")
	announcementId, _ := strconv.Atoi(announcementIdStr)

	if err := con.groupService.DeleteAnnouncement(uint(groupId), uint(announcementId)); err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c)
}

// GetAnnouncement 获取群组公告
// @Summary 查询群组公告
// @Description 获取指定群组的公告内容
// @Tags Group
// @Produce json
// @Param group_id path int true "群组ID"
// @Success 200 {object} model.Response{data=response.GroupAnnouncement}
// @Router /group/{group_id}/announcement [get]
func (con GroupController) GetAnnouncement(c *gin.Context) {
	groupIdStr := c.Param("group_id")
	groupId, _ := strconv.Atoi(groupIdStr)

	announcement, err := con.groupService.GetAnnouncement(uint(groupId))
	if err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c, announcement)
}

// GetAnnouncementList 获取群组公告列表
// @Summary 获取群组公告列表
// @Description 获取指定群组的所有公告
// @Tags Group
// @Produce json
// @Param group_id path int true "群组ID"
// @Success 200 {object} model.Response{data=[]response.GroupAnnouncement}
// @Router /group/{group_id}/announcement_list [get]
func (con GroupController) GetAnnouncementList(c *gin.Context) {
	groupIdStr := c.Param("group_id")
	groupId, _ := strconv.Atoi(groupIdStr)

	announcements, err := con.groupService.GetAnnouncementList(uint(groupId))
	if err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c, announcements)
}

func (con GroupController) KickMember(c *gin.Context) {
	var req request.KickMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}

	userId := c.GetUint("id")

	if err := con.groupService.KickMember(userId, req.GroupID, req.MemberID); err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c)
}

func (con GroupController) SetAdmin(c *gin.Context) {
	var req request.SetAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}

	userId := c.GetUint("id")

	if err := con.groupService.SetAdmin(userId, req.GroupID, req.MemberID); err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c)
}

func (con GroupController) UnsetAdmin(c *gin.Context) {
	var req request.UnsetAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}

	operatorId := c.GetUint("id")
	groupIdStr := c.Param("group_id")
	groupId, err := strconv.ParseUint(groupIdStr, 10, 64)
	if err != nil {
		con.Error(c, "无效的 group_id")
		return
	}

	if err := con.groupService.UnsetAdmin(operatorId, uint(groupId), req.MemberID); err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c)
}

func (con GroupController) MuteMember(c *gin.Context) {
	var req request.MuteMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}

	operatorId := c.GetUint("id")
	groupIdStr := c.Param("group_id")
	groupId, err := strconv.ParseUint(groupIdStr, 10, 64)
	if err != nil {
		con.Error(c, "无效的 group_id")
		return
	}

	if err := con.groupService.MuteMember(operatorId, uint(groupId), req.MemberID, req.Duration); err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c)
}

func (con GroupController) UnmuteMember(c *gin.Context) {
	var req request.UnmuteMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		con.Error(c, err.Error())
		return
	}
	userId := c.GetUint("id")
	groupId := c.GetUint("group_id")

	if err := con.groupService.UnmuteMember(userId, groupId, req.MemberID); err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c)
}
