package user

import (
	"context"
	"errors"

	uuid "github.com/nu7hatch/gouuid"
	pb "github.com/u03013112/ss-pb/user"
)

// Srv ：服务
type Srv struct{}

// Login :
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

// GetRoles :
func (s *Srv) GetRoles(ctx context.Context, in *pb.GetRolesRequest) (*pb.GetRolesReply, error) {
	u, err := getUserInfoByToken(in.Token)
	if err != nil {
		return nil, err
	}

	return &pb.GetRolesReply{Role: u.Role}, nil
}

// GetUserInfo :
func (s *Srv) GetUserInfo(ctx context.Context, in *pb.GetUserInfoRequest) (*pb.GetUserInfoReply, error) {
	// 各种用户都查一下，目前就只有android，就只能查android
	str, err := grpcGetAndroidUserStatus(in.Token)
	if err == nil {
		ret := &pb.GetUserInfoReply{
			Type:   "android",
			Status: str,
		}
		return ret, nil
	}
	return nil, errors.New("not found user")
}
