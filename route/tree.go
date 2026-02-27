package route

import (
	"strings"

	"github.com/go-web/public"
)

type Tree struct {
	root *Node
}

func NewTree() *Tree {
	return &Tree{root: &Node{children: make(map[string]*Node)}}
}

// AddRoute 注册路由
func (r *Tree) AddNode(path string, handler public.HandlerFunc) {
	parts := parsePath(path)
	curr := r.root
	for _, part := range parts {
		if curr.children[part] == nil {
			isParam := strings.HasPrefix(part, ":")
			curr.children[part] = &Node{
				path:     part,
				children: make(map[string]*Node),
				isParam:  isParam,
			}
		}
		curr = curr.children[part]
	}
	curr.isEnd = true
	curr.handler = handler
	curr.fullPath = path
}

// Search 搜索路径是否存在
func (r *Tree) GetNode(path string) *Node {
	parts := parsePath(path)
	curr := r.root

	for _, part := range parts {
		// 1. 精确匹配
		if next, ok := curr.children[part]; ok {
			curr = next
			continue
		}

		// 2. 参数匹配逻辑 (查找是否存在以 : 开头的子节点)
		foundParam := false
		for _, child := range curr.children {
			if child.isParam {
				curr = child
				foundParam = true
				break
			}
		}

		if !foundParam {
			return nil
		}
	}
	if curr.isEnd {
		return curr
	} else {
		return nil
	}
}

func parsePath(path string) []string {
	res := make([]string, 0)
	for _, s := range strings.Split(path, "/") {
		if s != "" {
			res = append(res, s)
		}
	}
	return res
}
