package route

import "github.com/cyuzhe1994-commits/go-web/public"

const name = ""

type Node struct {
	path     string             // 当前节点的路径片段，例如 "user" 或 ":id"
	children map[string]*Node   // 子节点，key 是路径片段
	isEnd    bool               // 是否是一个完整的注册路径
	handler  public.HandlerFunc // 到达该节点时执行的业务逻辑
	isParam  bool               // 是否是参数节点（如 :id）
	fullPath string             // 完整的注册路径，例如 "/user/:id"
}

func (n *Node) GetHandler() public.HandlerFunc {
	if n.isEnd {
		return n.handler
	}
	return nil
}

func (n *Node) GetFullPath() string {
	if n.isEnd {
		return n.fullPath
	}
	return ""
}
