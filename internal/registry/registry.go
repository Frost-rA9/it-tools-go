// Package registry 定义工具的统一接口与注册表，聚合所有工具供前端发现与调用。
package registry

import "context"

// Tool 描述一个工具的可发现元数据。
type Tool struct {
	ID          string   // 唯一标识，如 "uuid-generator"
	Name        string   // 展示名称
	Description string   // 一句话描述
	Category    string   // 所属分类（见 SPEC.md §6）
	Keywords    []string // 搜索关键词
}

// Executor 定义工具的执行逻辑，input/output 由各工具自行约定。
type Executor interface {
	Execute(ctx context.Context, input any) (output any, err error)
}

type entry struct {
	meta     Tool
	executor Executor
}

// Registry 聚合所有已注册的工具。
type Registry struct {
	entries []entry
}

// New 创建一个空的注册表。
func New() *Registry {
	return &Registry{}
}

// Register 向注册表加入一个工具及其执行逻辑。
func (r *Registry) Register(t Tool, e Executor) {
	r.entries = append(r.entries, entry{meta: t, executor: e})
}

// List 返回全部工具的元数据，供前端渲染侧边栏。
func (r *Registry) List() []Tool {
	metas := make([]Tool, 0, len(r.entries))
	for _, e := range r.entries {
		metas = append(metas, e.meta)
	}
	return metas
}

// Get 按 ID 查找工具的元数据。
func (r *Registry) Get(id string) (Tool, bool) {
	for _, e := range r.entries {
		if e.meta.ID == id {
			return e.meta, true
		}
	}
	return Tool{}, false
}

// Execute 按 ID 执行工具，未找到时返回 ok=false。
func (r *Registry) Execute(ctx context.Context, id string, input any) (output any, ok bool, err error) {
	for _, e := range r.entries {
		if e.meta.ID == id {
			out, err := e.executor.Execute(ctx, input)
			return out, true, err
		}
	}
	return nil, false, nil
}
