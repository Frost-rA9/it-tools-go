// Package registry 定义工具的统一接口与注册表，聚合所有工具供前端发现与调用。
package registry

import "context"

// 工具分类（对齐 SPEC.md §6，中文展示）。
const (
	CategoryCrypto       = "加密"
	CategoryConverter    = "转换器"
	CategoryWeb          = "Web"
	CategoryImagesVideos = "图片和视频"
	CategoryDevelopment  = "开发"
	CategoryNetwork      = "网络"
	CategoryMath         = "数学"
	CategoryMeasurement  = "测量"
	CategoryText         = "文本"
	CategoryData         = "数据"
)

// Tool 描述一个工具的可发现元数据。
type Tool struct {
	ID          string   `json:"id"`          // 唯一标识，如 "base64-string-converter"
	Name        string   `json:"name"`        // 展示名称
	Description string   `json:"description"` // 一句话描述
	Category    string   `json:"category"`    // 所属分类
	Keywords    []string `json:"keywords"`    // 搜索关键词
}

// Executor 定义工具的执行逻辑。input 与 output 均为 JSON 字符串，结构由各工具自行约定。
type Executor interface {
	Execute(ctx context.Context, input string) (output string, err error)
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

// List 返回全部工具的元数据，供前端渲染侧边栏与首页。
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
func (r *Registry) Execute(ctx context.Context, id, input string) (output string, ok bool, err error) {
	for _, e := range r.entries {
		if e.meta.ID == id {
			out, err := e.executor.Execute(ctx, input)
			return out, true, err
		}
	}
	return "", false, nil
}
