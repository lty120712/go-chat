package service

import (
	"errors"
	"go-chat/internal/db"
	"go-chat/internal/model"
	request "go-chat/internal/model/request"
	response "go-chat/internal/model/response"
	"go-chat/internal/repository"
	"go-chat/internal/utils/idUtil"
	"gorm.io/gorm"
	"sort"
	"sync"
)

type GroupService struct{}

var (
	groupServiceInstance *GroupService
	groupOnce            sync.Once
)

func GetGroupService() *GroupService {
	groupOnce.Do(func() {
		groupServiceInstance = &GroupService{}
	})
	return groupServiceInstance
}

// Create 创建群组
func (s GroupService) Create(req *request.GroupCreateRequest) error {
	err := db.Mysql.Transaction(func(tx *gorm.DB) error {
		groupRepo := repository.GetGroupRepository()
		groupMemberRepo := repository.GetGroupMemberRepository()
		userRepo := repository.GetUserRepository()

		maxAttempts := 3
		code := idUtil.GenerateId()
		for attempts := 0; attempts < maxAttempts; attempts++ {
			if exists := groupRepo.ExistsByCode(code, tx); !exists {
				break
			}
			code = idUtil.GenerateId()
		}

		if exists := groupRepo.ExistsByCode(code, tx); exists {
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

		if err := groupRepo.Save(group, tx); err != nil {
			return err
		}

		if req.MemberList != nil && len(*req.MemberList) > 0 {
			groupMemberList := make([]*model.GroupMember, len(*req.MemberList))

			idNickNameMap, err := userRepo.GetNickNamesByIds(*req.MemberList, tx)
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

			if err := groupMemberRepo.SaveBatch(groupMemberList, tx); err != nil {
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
		userRepo := repository.GetUserRepository()
		groupMemberRepo := repository.GetGroupMemberRepository()
		exists := groupMemberRepo.ExistsByGroupIdAndUserId(groupId, userId, tx)
		if exists {
			return errors.New("用户已加入该群组")
		}
		rejoin := groupMemberRepo.RejoinGroupIfDeleted(groupId, userId, tx)
		if rejoin {
			return nil
		}
		nickname, err := userRepo.GetNickNamesById(userId, tx)
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
		if err := groupMemberRepo.Save(member, tx); err != nil {
			return err
		}
		return nil
	})
}

func (s GroupService) Quit(groupId uint, memberId uint) error {
	memberRepository := repository.GetGroupMemberRepository()
	//群主不能退
	if memberRepository.IsOwner(groupId, memberId) {
		return errors.New("owner can not quit")
	}
	return memberRepository.DeleteByGroupIdAndUserId(groupId, memberId)
}

func (s GroupService) Member(groupId uint) (memberList []response.MemberVo, err error) {
	memberRepository := repository.GetGroupMemberRepository()
	memberList, err = memberRepository.GetMemberListByGroupId(groupId)
	if err != nil {
		return
	}
	sort.Slice(memberList, func(i, j int) bool {
		return memberList[i].Role > memberList[j].Role
	})
	return
}
