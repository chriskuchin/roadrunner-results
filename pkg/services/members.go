package services

import (
	"context"

	"github.com/chriskuchin/roadrunner-results/pkg/db"
)

var membersInstance *MembersService

type MembersService struct {
	membersDao *db.MemberDao
}

func NewMemberService(dao *db.MemberDao) {
	if membersInstance == nil {
		membersInstance = &MembersService{
			membersDao: dao,
		}
	}
}

func GetMembersServiceInstance() *MembersService {
	return membersInstance
}

func (m *MembersService) UpsertMember(ctx context.Context, firstName, lastName, gender string, birthYear int) error {
	return m.membersDao.UpsertMember(ctx, firstName, lastName, gender, birthYear)
}
