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

// init 是 Go 语言特殊的包初始化函数，在包被导入时自动执行。
// Go 的 init 机制：每个包可以有多个 init()，按照编译顺序在 main() 之前执行。
// 这里利用 init() 在程序启动时自动将 User 业务逻辑实现注册到全局的 service 层，
// 这样 controller 层只需要调用 service.User() 即可获取实现，无需手动创建依赖。
func init() {
	service.RegisterUser(NewUser())
}

// NewUser 是构造函数，返回 IUser 接口的实现。
// Go 语言约定：NewXxx 是创建结构体实例的常见命名模式。
// 返回接口类型 service.IUser 而非具体类型，实现依赖倒置（调用方依赖接口而非实现）。
func NewUser() service.IUser {
	return &sUser{}
}

// sUser 是 User 业务逻辑的实现结构体。
// 注意首字母小写：Go 中首字母小写的类型/函数为包内私有（unexported），
// 外部只能通过接口调用，强制了封装——调用方无法直接创建或访问这个结构体。
type sUser struct{}

// Create 创建新用户。
// 流程：检查用户名唯一性 → 加密密码 → 插入数据库。
// dao.User.Ctx(ctx) 是 GoFrame ORM 的链式调用起始：
//   Ctx() 绑定上下文（用于链路追踪/超时控制），
//   Where() 添加查询条件，
//   Count() 返回匹配行数。
// goframe 链式调用中每个方法都返回 Model 对象，可以继续链式调用。
func (s *sUser) Create(ctx context.Context, req *v1.UserCreateReq) (*v1.UserCreateRes, error) {
	count, err := dao.User.Ctx(ctx).Where(dao.User.Columns.Username, req.Username).Count()
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

	// dao.User.Data(&do.User{...}).Insert() 是 GoFrame ORM 插入操作。
	// Data() 接收一个结构体或 map 作为要插入的数据，
	// Insert() 执行 INSERT SQL 并返回结果。
	// do.User 是 GoFrame 的"数据操作"结构体，与数据库字段一一对应。
	// gtime.Now() 用于设置创建/更新时间，gtime 是 GoFrame 提供的增强版 time.Time。
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

// Delete 执行"软删除"——并非真正从数据库删除记录，而是设置 DeletedAt 字段为当前时间。
// 软删除的好处：数据可恢复、保留历史记录、不影响关联数据的完整性。
// 所有后续查询都通过 WhereNull(DeletedAt) 过滤掉已删除的记录。
// gtime.Now() 生成删除时间戳，用于标识记录被删除的时刻。
func (s *sUser) Delete(ctx context.Context, req *v1.UserDeleteReq) error {
	_, err := dao.User.Ctx(ctx).Where(dao.User.Columns.Id, req.Id).Update(gdb.Map{
		dao.User.Columns.DeletedAt: gtime.Now(),
	})
	return err
}

// Update 更新用户信息。只更新请求中提供的字段（非空字符串才覆盖），
// 然后统一设置 UpdatedAt 为当前时间，标记最后修改时间。
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

	_, err := dao.User.Ctx(ctx).Where(dao.User.Columns.Id, req.Id).Data(data).Update()
	return err
}

// GetOne 查询单个用户。
// dao.User.Ctx(ctx).Where(...).WhereNull(...).Scan(&user) 是 GoFrame ORM 链式查询：
//   Ctx(ctx)     → 绑定上下文
//   Where(列, 值)  → WHERE 条件
//   WhereNull(列)  → WHERE 列 IS NULL（过滤已软删除记录）
//   Scan(&user)   → 将查询结果扫描到 struct 指针
// Scan 是 GoFrame 的"瑞士军刀"——自动匹配字段名到结构体字段，
// 支持单条记录（结构体指针）和集合（切片指针）。
func (s *sUser) GetOne(ctx context.Context, req *v1.UserGetOneReq) (*v1.UserGetOneRes, error) {
	var user *entity.User
	err := dao.User.Ctx(ctx).Where(dao.User.Columns.Id, req.Id).WhereNull(dao.User.Columns.DeletedAt).Scan(&user)
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

// GetRoles 查询指定用户已分配的角色 ID 列表。
// 从多对多关联表 user_role 中查询，只取 RoleId 字段。
// Fields() 指定要查询的列，避免 SELECT * 带来的性能开销。
func (s *sUser) GetRoles(ctx context.Context, req *v1.UserGetRolesReq) (*v1.UserGetRolesRes, error) {
	var roleIds []uint64
	err := dao.UserRole.Ctx(ctx).
		Fields(dao.UserRole.Columns.RoleId).
		Where(dao.UserRole.Columns.UserId, req.Id).
		Scan(&roleIds)
	if err != nil {
		return nil, err
	}
	return &v1.UserGetRolesRes{RoleIds: roleIds}, nil
}

// AssignRoles 分配角色给用户。
// 使用 Transaction（数据库事务）保证原子性：先删除该用户所有旧角色，再插入新角色。
// 如果中间失败，事务自动回滚，不会出现用户角色数据不一致的情况。
// tx.Model(...) 用法与 dao.UserRole.Ctx(ctx) 等效，但绑定了事务对象。
func (s *sUser) AssignRoles(ctx context.Context, req *v1.UserAssignRolesReq) error {
	return dao.UserRole.Ctx(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := tx.Model(dao.UserRole.Table).Where(dao.UserRole.Columns.UserId, req.Id).Delete()
		if err != nil {
			return err
		}

		if len(req.RoleIds) == 0 {
			return nil
		}

		data := make([]gdb.Map, 0, len(req.RoleIds))
		for _, roleId := range req.RoleIds {
			data = append(data, gdb.Map{
				dao.UserRole.Columns.UserId: req.Id,
				dao.UserRole.Columns.RoleId: roleId,
			})
		}

		_, err = tx.Model(dao.UserRole.Table).Data(data).Insert()
		return err
	})
}

// List 分页查询用户列表，支持按用户名模糊搜索和按状态筛选。
//
// Count / Page / OrderDesc / Scan 模式详解：
// 1. 先调用 Count() 获取总记录数（用于前端分页组件显示总页数）。
//    注意：Count() 会"消费"当前 Model，之后必须重新从 m 开始链式调用。
//    这是因为 GoFrame 的 Count() 内部执行 SELECT COUNT(*) 后，Model 状态已改变。
// 2. 再通过 Page(page, pageSize) 设置 LIMIT/OFFSET 进行分页。
// 3. OrderDesc(列) 按指定列降序排列（最新数据在前）。
// 4. Scan(&list) 将结果扫描到切片指针。
//
// req.Status 的类型是 *int（指针），这是为了区分"传了 0"和"没传"。
// 如果使用 uint，Go 零值为 0，无法表示"未指定"状态，导致无法过滤禁用用户。
func (s *sUser) List(ctx context.Context, req *v1.UserListReq) (*v1.UserListRes, error) {
	m := dao.User.Ctx(ctx).WhereNull(dao.User.Columns.DeletedAt)

	if req.Username != "" {
		m = m.WhereLike(dao.User.Columns.Username, "%"+req.Username+"%")
	}
	if req.Status != nil {
		m = m.Where(dao.User.Columns.Status, *req.Status)
	}

	var total int
	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	var list []*entity.User
	err = m.Page(req.Page, req.PageSize).OrderDesc(dao.User.Columns.Id).Scan(&list)
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
