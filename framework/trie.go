package framework

import (
	"errors"
	"fmt"
	"strings"
)

type Tree struct {
	root *node                 // 根节点
}

type node struct {
	isLast bool                 //该节点是否能成为一个独立的uri, 是否自身就是一个终极节点
	segment string             // uri中的字符串
	handler ControllerHandler  // 控制器
	childs  []*node            // 子节点
}

func newNode()*node  {
	return &node{
		isLast:false,
		segment:"",
		childs:[]*node{},
	}
}

func NewTree() *Tree {
	root := newNode()
	return &Tree{root}
}

// 判断一个segment是否是通用segment，即以:开头
func isWildSegment(segment string)bool  {
	return strings.HasPrefix(segment,":")
}

// 过滤下一层满足segment规则的子节点
func (n *node)filterChildNodes(segment string)[]*node  {
	if len(n.childs) == 0{
		return nil
	}

	//如果segment是通配符，则所有下一层子节点都满足需求
	if isWildSegment(segment){
		return n.childs
	}

	nodes := make([]*node,0,len(n.childs))

	// 过滤所有的下一层子节点
	for _,cnode := range n.childs{
		if isWildSegment(cnode.segment){
			// 如果下一层子节点有通配符，则满足需求
			nodes = append(nodes,cnode)
		}else if cnode.segment == segment{
			// 如果下一层子节点没有通配符，但是文本完全匹配，则满足需求
			nodes = append(nodes, cnode)
		}
	}
	return nodes
}

// 判断路由是否已经在节点的所有子节点树中存在了
func(n *node)matchNode(uri string)*node {
	// 使用分隔符将uri切割为两个部分
	fmt.Println("uri:",uri)
	segments := strings.SplitN(uri,"/",2)
	fmt.Println("segments",segments)

	// 第一个部分用于匹配下一层子节点
	segment := segments[0]
	fmt.Println("segment:  ",segment)

	if !isWildSegment(segment){
		segment = strings.ToUpper(segment)
	}
	fmt.Println("isWildSegment",isWildSegment(segment))


	// 匹配符合的下一层子节点
	cnodes := n.filterChildNodes(segment)
	fmt.Println("cnodes",cnodes)
	// 如果当前子节点没有一个符合，那么说明这个uri一定是之前不存在, 直接返回nil
	if cnodes == nil || len(cnodes) == 0{
		return nil
	}

	//如果只有一个segment，则是最后一个标记
	if len(segments) == 1{
		// 如果segment已经是最后一个节点，判断这些cnode是否有isLast标志
		for _,tn := range cnodes{
			if tn.isLast{
				return tn
			}
		}

		//都不是最后一个节点
		return nil
	}

	//如果有2个segment,递归每个子节点继续进行查找
	for _, tn := range cnodes {
		tnMatch := tn.matchNode(segments[1])
		if tnMatch != nil {
			return tnMatch
		}
	}
	return nil
}

//  增加路由节点
func (tree *Tree)AddRouter(uri string,handler ControllerHandler)error  {
	n := tree.root
	fmt.Println("add router")
	fmt.Println("-----",uri)
	//确认路由是否冲突
	if n.matchNode(uri) != nil{
		return errors.New("route exist: "+uri)
	}

	segments := strings.Split(uri, "/")
	fmt.Println("segments:",segments)

	// 对每个segment
	for index,segment := range segments{
		fmt.Println("index:",index)
		fmt.Println("segment:",segment)

		// 最终进入Node segment的字段
		if !isWildSegment(segment){
			fmt.Println("isWildSegment")
			segment = strings.ToUpper(segment)
			fmt.Println(segment)
		}
		fmt.Println("len(segments)",len(segments))

		isLast := index == len(segments) -1  //最后一个为ture

		var objNode *node // 标记是否有合适的子节点

		childNodes := n.filterChildNodes(segment)

		// 如果有匹配的子节点
		if len(childNodes) > 0 {
			// 如果有segment相同的子节点，则选择这个子节点
			for _, cnode := range childNodes {
				if cnode.segment == segment {
					objNode = cnode
					break
				}
			}
		}

		if objNode == nil {
			// 创建一个当前node的节点
			cnode := newNode()
			cnode.segment = segment
			if isLast {
				cnode.isLast = true
				cnode.handler = handler
				//cnode.segment = "我是第一个"
			}

			n.childs = append(n.childs, cnode)
			objNode = cnode
		}

		fmt.Println("objNode 赋值:",objNode)
		n = objNode
	}

	//fmt.Println("n:",n)
	//fmt.Println("end")


	//fmt.Println("tree",tree.root.childs)
	//for _,nodeItem := range tree.root.childs{
	//	fmt.Println("nodeItem:",nodeItem.segment)
	//}
	//os.Exit(3)
	return nil
}

// 匹配uri
func (tree *Tree) FindHandler(uri string) ControllerHandler {
	matchNode := tree.root.matchNode(uri)
	fmt.Println("matchNode :",matchNode)
	if matchNode == nil {
		return nil
	}
	return matchNode.handler
}























