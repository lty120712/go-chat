package service

import (
	"errors"
	"go-chat/internal/db"
	interfacerepository "go-chat/internal/interfaces/repository"
	"go-chat/internal/model"
	request "go-chat/internal/model/request"
	response "go-chat/internal/model/response"
	"go-chat/internal/utils/idUtil"
	"gorm.io/gorm"
	"sort"
	"sync"
)

type GroupService struct {
	groupRepository       interfacerepository.GroupRepositoryInterface
	messageRepository     interfacerepository.MessageRepositoryInterface
	userRepository        interfacerepository.UserRepositoryInterface
	groupMemberRepository interfacerepository.GroupMemberRepositoryInterface
}

var (
	GroupServiceInstance *GroupService
	groupOnce            sync.Once
)

func InitGroupService(groupRepository interfacerepository.GroupRepositoryInterface,
	messageRepository interfacerepository.MessageRepositoryInterface,
	userRepository interfacerepository.UserRepositoryInterface,
	groupMemberRepository interfacerepository.GroupMemberRepositoryInterface) {
	groupOnce.Do(func() {
		GroupServiceInstance = &GroupService{
			groupRepository:       groupRepository,
			messageRepository:     messageRepository,
			userRepository:        userRepository,
			groupMemberRepository: groupMemberRepository,
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
					Status:    model.Enable,
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
			Status:    model.Enable,
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
