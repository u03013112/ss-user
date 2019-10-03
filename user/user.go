package user

import (
	"context"
	"errors"

	uuid "github.com/nu7hatch/gouuid"
	pb "github.com/u03013112/ss-pb/user"
)

// Srv ：服务
type Srv struct{}

func (s *Srv) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginReply, error) {
	id, err := auth(in.Username, in.Passwd)
	if err != nil {
		return &pb.LoginReply{}, errors.New("auth failed")
	}

	token, _ := uuid.NewV4()
	u := getUserInfo(id)
	u.Token = token.String()
	updateUserInfo(u)
	return &pb.LoginReply{Token: token.String()}, nil
}

func (s *Srv) GetRoles(ctx context.Context, in *pb.GetRolesRequest) (*pb.GetRolesReply, error) {
	u, err := getUserInfoByToken(in.Token)
	if err != nil {
		return nil, err
	}

	return &pb.GetRolesReply{Role: u.Role}, nil
}
