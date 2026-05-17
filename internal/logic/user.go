package logic

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	v1 "tool-go/api/v1"
	"tool-go/internal/dao"
	"tool-go/internal/library/password"
	"tool-go/internal/model/do"
	"tool-go/internal/model/entity"
	"tool-go/internal/service"
)

func init() {
	service.RegisterUser(NewUser())
}

func NewUser() service.IUser {
	return &sUser{}
}

type sUser struct{}

func (s *sUser) Create(ctx context.Context, req *v1.UserCreateReq) (*v1.UserCreateRes, error) {
	var count int
	err := dao.User.Ctx(ctx).Where(dao.User.columns.Username, req.Username).Count(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, gerror.New("用户名已存在")
	}

	hash, salt, err := password.CreatePassword(req.Password)
	if err != nil {
		return nil, gerror.New("密码加密失败")
	}

	result, err := dao.User.Data(&do.User{
		Username:  req.Username,
		Password:  hash,
		Salt:      salt,
		Nickname:  req.Nickname,
		Email:     req.Email,
		Phone:     req.Phone,
		Status:    req.Status,
		CreatedAt: gtime.Now(),
		UpdatedAt: gtime.Now(),
	}).Insert()
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return &v1.UserCreateRes{Id: uint64(id)}, nil
}

func (s *sUser) Delete(ctx context.Context, req *v1.UserDeleteReq) error {
	_, err := dao.User.Ctx(ctx).Where(dao.User.columns.Id, req.Id).Update(gdb.Map{
		dao.User.columns.DeletedAt: gtime.Now(),
	})
	return err
}

func (s *sUser) Update(ctx context.Context, req *v1.UserUpdateReq) error {
	data := do.User{}
	if req.Username != "" {
		data.Username = req.Username
	}
	if req.Nickname != "" {
		data.Nickname = req.Nickname
	}
	if req.Email != "" {
		data.Email = req.Email
	}
	if req.Phone != "" {
		data.Phone = req.Phone
	}
	data.Status = req.Status
	data.UpdatedAt = gtime.Now()

	_, err := dao.User.Ctx(ctx).Where(dao.User.columns.Id, req.Id).Data(data).Update()
	return err
}

func (s *sUser) GetOne(ctx context.Context, req *v1.UserGetOneReq) (*v1.UserGetOneRes, error) {
	var user *entity.User
	err := dao.User.Ctx(ctx).Where(dao.User.columns.Id, req.Id).WhereNull(dao.User.columns.DeletedAt).Scan(&user)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, gerror.New("用户不存在")
	}

	return &v1.UserGetOneRes{
		Id:        user.Id,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Phone:     user.Phone,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}, nil
}

func (s *sUser) List(ctx context.Context, req *v1.UserListReq) (*v1.UserListRes, error) {
	m := dao.User.Ctx(ctx).WhereNull(dao.User.columns.DeletedAt)

	if req.Username != "" {
		m = m.WhereLike(dao.User.columns.Username, "%"+req.Username+"%")
	}
	if req.Status > 0 {
		m = m.Where(dao.User.columns.Status, req.Status)
	}

	var total int
	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	var list []*entity.User
	err = m.Page(req.Page, req.PageSize).OrderDesc(dao.User.columns.Id).Scan(&list)
	if err != nil {
		return nil, err
	}

	items := make([]v1.UserItem, 0, len(list))
	for _, u := range list {
		items = append(items, v1.UserItem{
			Id:        u.Id,
			Username:  u.Username,
			Nickname:  u.Nickname,
			Email:     u.Email,
			Phone:     u.Phone,
			Status:    u.Status,
			CreatedAt: u.CreatedAt.String(),
		})
	}

	return &v1.UserListRes{
		List:  items,
		Total: total,
		Page:  req.Page,
	}, nil
}
