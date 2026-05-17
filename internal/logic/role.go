package logic

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	v1 "tool-go/api/v1"
	"tool-go/internal/dao"
	"tool-go/internal/model/do"
	"tool-go/internal/model/entity"
	"tool-go/internal/service"
)

// init 是 Go 语言的包初始化函数，在包被导入时自动执行。
// Go 的 init 机制：每个包可以有多个 init()，按照编译顺序在 main() 之前执行。
// 这里在程序启动时自动将 Role 业务逻辑实现注册到全局的 service 层，
// controller 只需调用 service.Role() 即可获取实现，无需手工管理依赖。
func init() {
	service.RegisterRole(NewRole())
}

// NewRole 是构造函数，返回 IRole 接口的实现。
// Go 语言约定：NewXxx 是创建结构体实例的常见命名模式。
// 返回接口类型 service.IRole 而非具体类型 sRole，实现依赖倒置。
func NewRole() service.IRole {
	return &sRole{}
}

// sRole 是 Role 业务逻辑的实现结构体。
// 首字母小写（unexported），外部包无法直接访问，
// 只能通过 service.Role() 获取接口实例，强制封装。
type sRole struct{}

// Create 创建新角色。
// GoFrame ORM 链式调用：dao.Role.Ctx(ctx) → Where() 条件判断 → Count() 查重，
// 然后通过 Data(&do.Role{}).Insert() 写入数据库。
// do.Role 是 GoFrame 的数据操作结构体，字段对应数据库 role 表。
func (s *sRole) Create(ctx context.Context, req *v1.RoleCreateReq) (*v1.RoleCreateRes, error) {
	count, err := dao.Role.Ctx(ctx).Where(dao.Role.Columns.Code, req.Code).Count()
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, gerror.New("角色编码已存在")
	}

	// Data() 接收数据操作结构体，Insert() 执行 INSERT SQL。
	// gtime.Now() 设置创建时间和更新时间，由 GoFrame 的 gtime 包提供。
	result, err := dao.Role.Data(&do.Role{
		Name:      req.Name,
		Code:      req.Code,
		Sort:      req.Sort,
		Status:    req.Status,
		Desc:      req.Desc,
		CreatedAt: gtime.Now(),
		UpdatedAt: gtime.Now(),
	}).Insert()
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return &v1.RoleCreateRes{Id: uint64(id)}, nil
}

// Delete 执行软删除——设置 DeletedAt 字段为当前时间，而非真正删除。
// 软删除后的记录通过 WhereNull(DeletedAt) 过滤，对业务透明。
func (s *sRole) Delete(ctx context.Context, req *v1.RoleDeleteReq) error {
	_, err := dao.Role.Ctx(ctx).Where(dao.Role.Columns.Id, req.Id).Update(gdb.Map{
		dao.Role.Columns.DeletedAt: gtime.Now(),
	})
	return err
}

// Update 更新角色信息。只覆盖请求中非空的字段，并刷新 UpdatedAt 时间戳。
func (s *sRole) Update(ctx context.Context, req *v1.RoleUpdateReq) error {
	data := do.Role{}
	if req.Name != "" {
		data.Name = req.Name
	}
	if req.Code != "" {
		data.Code = req.Code
	}
	data.Sort = req.Sort
	data.Status = req.Status
	data.Desc = req.Desc
	data.UpdatedAt = gtime.Now()

	_, err := dao.Role.Ctx(ctx).Where(dao.Role.Columns.Id, req.Id).Data(data).Update()
	return err
}

// GetOne 查询单个角色。
// Scan(&role) 将 SQL 查询结果自动映射到 entity.Role 结构体。
// 如果记录不存在或已软删除，role 为 nil，返回错误提示。
func (s *sRole) GetOne(ctx context.Context, req *v1.RoleGetOneReq) (*v1.RoleGetOneRes, error) {
	var role *entity.Role
	err := dao.Role.Ctx(ctx).Where(dao.Role.Columns.Id, req.Id).WhereNull(dao.Role.Columns.DeletedAt).Scan(&role)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, gerror.New("role not found")
	}

	return &v1.RoleGetOneRes{
		Id:        role.Id,
		Name:      role.Name,
		Code:      role.Code,
		Sort:      role.Sort,
		Status:    role.Status,
		Desc:      role.Desc,
		CreatedAt: role.CreatedAt.String(),
		UpdatedAt: role.UpdatedAt.String(),
	}, nil
}

// List 分页查询角色列表，支持按名称模糊搜索和按状态筛选。
//
// Count / Page / Order / Scan 模式：
// 1. 必须先调用 Count() 获取总条数（Count 会消费 Model，之后需从 m 重新链式调用）。
// 2. 然后 Page() 分页、OrderAsc() / OrderDesc() 排序。
// 3. 最后 Scan() 扫描结果到切片。
// OrderAsc(dao.Role.Columns.Sort) 先按排序字段升序，再按 ID 降序，确保同级角色有序。
func (s *sRole) List(ctx context.Context, req *v1.RoleListReq) (*v1.RoleListRes, error) {
	m := dao.Role.Ctx(ctx).WhereNull(dao.Role.Columns.DeletedAt)

	if req.Name != "" {
		m = m.WhereLike(dao.Role.Columns.Name, "%"+req.Name+"%")
	}
	if req.Status != nil {
		m = m.Where(dao.Role.Columns.Status, *req.Status)
	}

	var total int
	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	var list []*entity.Role
	err = m.Page(req.Page, req.PageSize).OrderAsc(dao.Role.Columns.Sort).OrderDesc(dao.Role.Columns.Id).Scan(&list)
	if err != nil {
		return nil, err
	}

	items := make([]v1.RoleItem, 0, len(list))
	for _, r := range list {
		items = append(items, v1.RoleItem{
			Id:        r.Id,
			Name:      r.Name,
			Code:      r.Code,
			Sort:      r.Sort,
			Status:    r.Status,
			Desc:      r.Desc,
			CreatedAt: r.CreatedAt.String(),
		})
	}

	return &v1.RoleListRes{
		List:  items,
		Total: total,
		Page:  req.Page,
	}, nil
}
