package controller

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	v1 "tool-go/api/v1"
	"tool-go/internal/dao"
	"tool-go/internal/library/jwt"
	"tool-go/internal/middleware"
	"tool-go/internal/model/entity"
)

var Auth = cAuth{}

type cAuth struct{}

func (c *cAuth) Login(ctx context.Context, req *v1.LoginReq) (*v1.LoginRes, error) {
	var user *entity.User
	err := dao.User.Ctx(ctx).
		Where(dao.User.columns.Username, req.Username).
		WhereNull(dao.User.columns.DeletedAt).
		Scan(&user)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, gerror.New("用户名或密码错误")
	}

	if user.Password != req.Password {
		return nil, gerror.New("用户名或密码错误")
	}

	if user.Status == 0 {
		return nil, gerror.New("账号已被禁用")
	}

	var roles []string
	err = dao.UserRole.Ctx(ctx).
		Fields(dao.Role.columns.Code).
		Where(dao.UserRole.columns.UserId, user.Id).
		LeftJoin("role", dao.Role.table+".id="+dao.UserRole.table+"."+dao.UserRole.columns.RoleId).
		WhereNull(dao.Role.columns.DeletedAt).
		Array(&roles)
	if err != nil {
		g.Log().Error(ctx, "获取用户角色失败:", err)
	}
	if len(roles) == 0 {
		roles = []string{"user"}
	}

	jwtConfig := g.Cfg().MustGet(ctx, "jwt").MapStrVar()
	secret := jwtConfig["secret"].String()
	if secret == "" {
		secret = "tool-go-jwt-secret-key-change-in-production"
	}
	expires := jwtConfig["expires"].Duration()
	if expires == 0 {
		expires = 24 * time.Hour
	}
	issuer := jwtConfig["issuer"].String()
	if issuer == "" {
		issuer = "tool-go"
	}

	j := jwt.New(secret, expires, issuer)
	token, err := j.GenerateToken(user.Id, user.Username, roles)
	if err != nil {
		return nil, gerror.New("生成token失败")
	}

	return &v1.LoginRes{
		Token:    token,
		UserId:   user.Id,
		Username: user.Username,
		Nickname: user.Nickname,
		Roles:    roles,
	}, nil
}

func (c *cAuth) GetUserInfo(ctx context.Context, req *v1.GetUserInfoReq) (*v1.GetUserInfoRes, error) {
	userId := middleware.GetUserId(ctx)
	if userId == 0 {
		return nil, gerror.New("未登录")
	}

	var user *entity.User
	err := dao.User.Ctx(ctx).
		Where(dao.User.columns.Id, userId).
		WhereNull(dao.User.columns.DeletedAt).
		Scan(&user)
	if err != nil || user == nil {
		return nil, gerror.New("用户不存在")
	}

	var roles []string
	err = dao.UserRole.Ctx(ctx).
		Fields(dao.Role.columns.Code).
		Where(dao.UserRole.columns.UserId, user.Id).
		LeftJoin("role", dao.Role.table+".id="+dao.UserRole.table+"."+dao.UserRole.columns.RoleId).
		WhereNull(dao.Role.columns.DeletedAt).
		Array(&roles)
	if err != nil {
		g.Log().Error(ctx, "获取用户角色失败:", err)
	}
	if len(roles) == 0 {
		roles = []string{"user"}
	}

	return &v1.GetUserInfoRes{
		UserId:   user.Id,
		Username: user.Username,
		Nickname: user.Nickname,
		Roles:    roles,
	}, nil
}

func (c *cAuth) Logout(ctx context.Context, req *v1.LogoutReq) (*v1.LogoutRes, error) {
	return &v1.LogoutRes{}, nil
}
