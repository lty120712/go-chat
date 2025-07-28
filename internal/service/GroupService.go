package service

import (
	"errors"
	"fmt"
	"github.com/lty120712/gorm-pagination/pagination"
	"go-chat/internal/db"
	interfacerepository "go-chat/internal/interfaces/repository"
	"go-chat/internal/model"
	request "go-chat/internal/model/request"
	response "go-chat/internal/model/response"
	"go-chat/internal/repository"
	"go-chat/internal/utils/idUtil"
	"gorm.io/gorm"
	"sort"
	"sync"
	"time"
)

type GroupService struct {
	groupRepository             interfacerepository.GroupRepositoryInterface
	messageRepository           interfacerepository.MessageRepositoryInterface
	userRepository              interfacerepository.UserRepositoryInterface
	groupMemberRepository       interfacerepository.GroupMemberRepositoryInterface
	groupAnnouncementRepository interfacerepository.GroupAnnouncementRepositoryInterface
}

var (
	GroupServiceInstance *GroupService
	groupOnce            sync.Once
)

func InitGroupService(groupRepository interfacerepository.GroupRepositoryInterface,
	messageRepository interfacerepository.MessageRepositoryInterface,
	userRepository interfacerepository.UserRepositoryInterface,
	groupMemberRepository interfacerepository.GroupMemberRepositoryInterface,
	groupAnnouncementRepository interfacerepository.GroupAnnouncementRepositoryInterface,
) {
	groupOnce.Do(func() {
		GroupServiceInstance = &GroupService{
			groupRepository:             groupRepository,
			messageRepository:           messageRepository,
			userRepository:              userRepository,
			groupMemberRepository:       groupMemberRepository,
			groupAnnouncementRepository: groupAnnouncementRepository,
		}
	})
}

// Create 创建群组
func (s GroupService) Create(req *request.GroupCreateRequest) error {
	err := db.Mysql.Transaction(func(tx *gorm.DB) error {
		maxAttempts := 3
		code := idUtil.GenerateId()
		for attempts := 0; attempts < maxAttempts; attempts++ {
			if exists := s.groupRepository.ExistsByCode(code, tx); !exists {
				break
			}
			code = idUtil.GenerateId()
		}

		if exists := s.groupRepository.ExistsByCode(code, tx); exists {
			return errors.New("error code, please try again")
		}

		group := &model.Group{
			OwnerId: req.UserId,
			Name:    req.Name,
			MaxNum:  req.MaxNum,
			Code:    code,
			Desc:    "群主很懒,什么也没留下~",
			Status:  model.Enable,
		}

		if err := s.groupRepository.Save(group, tx); err != nil {
			return err
		}

		if req.MemberList != nil && len(*req.MemberList) > 0 {
			groupMemberList := make([]*model.GroupMember, len(*req.MemberList))

			idNickNameMap, err := s.userRepository.GetNickNamesByIds(*req.MemberList, tx)
			if err != nil {
				return err
			}

			for i, memberId := range *req.MemberList {
				nickname := idNickNameMap[memberId]
				var role model.Role
				if req.UserId == memberId {
					role = model.Owner // 群主角色
				} else {
					role = model.Member // 普通成员角色
				}
				groupMemberList[i] = &model.GroupMember{
					GroupId:   group.ID,
					MemberId:  memberId,
					GNickName: nickname,
					Role:      role,
				}
			}

			if err := s.groupMemberRepository.SaveBatch(groupMemberList, tx); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s GroupService) Join(groupId uint, userId uint) error {
	return db.Mysql.Transaction(func(tx *gorm.DB) error {
		exists := s.groupMemberRepository.ExistsByGroupIdAndUserId(groupId, userId, tx)
		if exists {
			return errors.New("用户已加入该群组")
		}
		rejoin := s.groupMemberRepository.RejoinGroupIfDeleted(groupId, userId, tx)
		if rejoin {
			return nil
		}
		nickname, err := s.userRepository.GetNickNamesById(userId, tx)
		if err != nil {
			return err
		}
		member := &model.GroupMember{
			GroupId:   groupId,
			MemberId:  userId,
			GNickName: nickname,
			Role:      model.Member,
		}
		if err := s.groupMemberRepository.Save(member, tx); err != nil {
			return err
		}
		return nil
	})
}

func (s GroupService) Quit(groupId uint, memberId uint) error {
	//群主不能退
	if s.groupMemberRepository.IsOwner(groupId, memberId) {
		return errors.New("owner can not quit")
	}
	return s.groupMemberRepository.DeleteByGroupIdAndUserId(groupId, memberId)
}

func (s GroupService) Member(groupId uint) (memberList []response.MemberVo, err error) {
	memberList, err = s.groupMemberRepository.GetMemberListByGroupId(groupId)
	if err != nil {
		return
	}
	sort.Slice(memberList, func(i, j int) bool {
		return memberList[i].Role > memberList[j].Role
	})
	return
}

// CreateAnnouncement 创建群组公告
func (s GroupService) CreateAnnouncement(groupId uint, req *request.GroupAnnouncementCreateRequest) error {

	// 插入公告数据
	announcement := &model.GroupAnnouncement{
		GroupID:   groupId,
		Content:   req.Content,
		Publisher: req.Publisher,
	}

	// 保存公告
	err := repository.GroupAnnouncementRepositoryInstance.Create(announcement)
	if err != nil {
		return fmt.Errorf("创建群组公告失败: %w", err)
	}

	return nil
}

// UpdateAnnouncement 更新群组公告
func (s GroupService) UpdateAnnouncement(groupId uint, req *request.GroupAnnouncementUpdateRequest) error {
	// 获取现有公告
	announcement, err := s.groupAnnouncementRepository.GetByID(uint(req.AnnouncementId))
	if err != nil {
		return fmt.Errorf("查询群组公告失败: %w", err)
	}
	if announcement.GroupID != groupId {
		return fmt.Errorf("群组ID不匹配")
	}

	// 更新公告内容
	announcement.Content = req.Content
	announcement.UpdatedAt = time.Now()

	// 保存更新后的公告
	err = repository.GroupAnnouncementRepositoryInstance.Update(announcement)
	if err != nil {
		return fmt.Errorf("更新群组公告失败: %w", err)
	}

	return nil
}

// DeleteAnnouncement 删除群组公告
func (s GroupService) DeleteAnnouncement(groupId uint, announcementId uint) error {
	// 获取现有公告
	announcement, err := s.groupAnnouncementRepository.GetByID(announcementId)
	if err != nil {
		return fmt.Errorf("查询群组公告失败: %w", err)
	}
	if announcement.GroupID != groupId {
		return fmt.Errorf("群组ID不匹配")
	}

	// 设置为删除状态
	err = repository.GroupAnnouncementRepositoryInstance.Delete(announcementId)
	if err != nil {
		return fmt.Errorf("删除群组公告失败: %w", err)
	}

	return nil
}

// GetAnnouncement 获取群组单个公告
func (s GroupService) GetAnnouncement(groupId uint) (*model.GroupAnnouncement, error) {
	// 查询群组的公告（可以按最新时间排序等）
	announcement, err := s.groupAnnouncementRepository.GetLatestByGroupID(groupId)
	if err != nil {
		return nil, fmt.Errorf("查询群组公告失败: %w", err)
	}

	return announcement, nil
}

// GetAnnouncementList 获取群组公告列表
func (s GroupService) GetAnnouncementList(groupId uint) ([]model.GroupAnnouncement, error) {
	// 查询群组公告列表
	announcements, err := s.groupAnnouncementRepository.GetListByGroupID(groupId)
	if err != nil {
		return nil, fmt.Errorf("获取群组公告列表失败: %w", err)
	}

	return announcements, nil
}

func (s GroupService) KickMember(operatorId, groupId, targetMemberId uint) error {
	// 获取操作者在该群的身份
	operator, err := s.groupMemberRepository.GetGroupMember(groupId, operatorId)
	if err != nil {
		return fmt.Errorf("获取操作人信息失败: %v", err)
	}
	if operator == nil {
		return errors.New("你不是该群成员")
	}
	if operator.Role != 1 && operator.Role != 2 {
		return errors.New("无权限踢人，仅群主或管理员可踢人")
	}

	// 群主只能被自己踢（禁止踢群主）
	target, err := s.groupMemberRepository.GetGroupMember(groupId, targetMemberId)
	if err != nil {
		return fmt.Errorf("获取目标成员失败: %v", err)
	}
	if target == nil {
		return errors.New("目标成员不存在")
	}
	if target.Role == 1 {
		return errors.New("不能踢群主")
	}
	if operator.Role == 2 && target.Role == 2 {
		return errors.New("管理员不能踢管理员")
	}

	// 删除成员记录（逻辑删除或物理删除皆可）
	if err := s.groupMemberRepository.RemoveMember(groupId, targetMemberId); err != nil {
		return fmt.Errorf("踢人失败: %v", err)
	}

	return nil
}

func (s GroupService) SetAdmin(operatorId, groupId, memberId uint) error {
	// 获取操作者信息
	operator, err := s.groupMemberRepository.GetGroupMember(groupId, operatorId)
	if err != nil {
		return fmt.Errorf("获取操作者信息失败: %v", err)
	}
	if operator == nil {
		return errors.New("操作者不是该群成员")
	}
	if operator.Role != 1 { // 只有群主有权限设置管理员
		return errors.New("只有群主才能设置管理员")
	}

	target, err := s.groupMemberRepository.GetGroupMember(groupId, memberId)
	if err != nil {
		return fmt.Errorf("获取目标成员失败: %v", err)
	}
	if target == nil {
		return errors.New("目标成员不存在")
	}
	if target.Role == model.Admin {
		return errors.New("群主不能被设置为管理员")
	}

	// 设置角色为管理员 (2)
	target.Role = 2
	err = s.groupRepository.UpdateRole(target)
	if err != nil {
		return fmt.Errorf("设置管理员失败: %v", err)
	}

	return nil
}

func (s GroupService) UnsetAdmin(operatorId, groupId, targetMemberId uint) error {
	operator, err := s.groupMemberRepository.GetGroupMember(groupId, operatorId)
	if err != nil {
		return fmt.Errorf("获取操作人信息失败: %v", err)
	}
	if operator == nil || operator.Role != 1 {
		return errors.New("仅群主可以取消管理员权限")
	}

	target, err := s.groupMemberRepository.GetGroupMember(groupId, targetMemberId)
	if err != nil {
		return fmt.Errorf("获取目标成员失败: %v", err)
	}
	if target == nil {
		return errors.New("目标成员不存在")
	}
	if target.Role != model.Admin {
		return errors.New("该成员不是管理员")
	}

	target.Role = model.Member // 设置为普通成员
	return s.groupRepository.UpdateRole(target)
}

func (s GroupService) MuteMember(operatorId, groupId, targetMemberId uint, duration int64) error {
	operator, err := s.groupMemberRepository.GetGroupMember(groupId, operatorId)
	if err != nil {
		return fmt.Errorf("获取操作人信息失败: %v", err)
	}
	if operator == nil || (operator.Role != model.Owner && operator.Role != model.Admin) {
		return errors.New("仅群主或管理员可进行禁言操作")
	}

	target, err := s.groupMemberRepository.GetGroupMember(groupId, targetMemberId)
	if err != nil {
		return fmt.Errorf("获取目标成员失败: %v", err)
	}
	if target == nil {
		return errors.New("目标成员不存在")
	}
	if target.Role == 1 {
		return errors.New("不能禁言群主")
	}
	if operator.Role == 2 && target.Role == 2 {
		return errors.New("管理员不能禁言管理员")
	}

	muteUntil := time.Now().Add(time.Duration(duration) * time.Second)
	target.MuteEnd = &muteUntil

	return s.groupRepository.UpdateRole(target)
}

func (s GroupService) UnmuteMember(operatorId, groupId, targetId uint) error {
	operator, err := s.groupMemberRepository.GetGroupMember(groupId, operatorId)
	if err != nil {
		return fmt.Errorf("获取操作人信息失败: %v", err)
	}
	if operator == nil {
		return errors.New("你不是该群成员")
	}
	if operator.Role != 1 && operator.Role != 2 {
		return errors.New("无权限操作，仅群主或管理员可解除禁言")
	}

	target, err := s.groupMemberRepository.GetGroupMember(groupId, targetId)
	if err != nil {
		return fmt.Errorf("获取目标成员失败: %v", err)
	}
	if target == nil {
		return errors.New("目标成员不存在")
	}
	if target.IsOwner() {
		return errors.New("无法解除群主禁言（群主默认不禁言）")
	}
	if operator.IsAdmin() && target.IsAdmin() {
		return errors.New("管理员不能解除管理员禁言")
	}

	target.MuteEnd = nil

	return s.groupRepository.UpdateRole(target)
}

func (s GroupService) Page(req request.GroupSearchRequest) (*pagination.PageResult[model.Group], error) {

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}
	res, err := s.groupRepository.Page(req)
	return res, err
}
