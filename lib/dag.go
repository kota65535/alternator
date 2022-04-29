package lib

import (
	"sort"
)

type Dag struct {
	Nodes     map[Alteration]bool
	NodeOrder LinkedList
	Deps      map[Alteration][]Alteration
}

func NewDag(_nodes []Alteration) Dag {
	// First, sort by SeqNum
	nodes := make([]Alteration, len(_nodes))
	copy(nodes, _nodes)
	sort.SliceStable(nodes, func(i, j int) bool { return nodes[i].SeqNum() < nodes[j].SeqNum() })
	nodeOrder := NewLinkedList(nodes)
	nodeMap := map[Alteration]bool{}
	for _, n := range nodes {
		nodeMap[n] = true
	}
	depMap := map[Alteration][]Alteration{}
	for _, n := range nodes {
		dependsOn := RemoveIf(n.DependsOn(), func(e Alteration) bool { return !nodeMap[e] })
		depMap[n] = dependsOn
	}
	return Dag{nodeMap, nodeOrder, depMap}
}

func (r Dag) Sort() []Alteration {
	var ret []Alteration
	for len(r.Nodes) > 0 {
		var head Alteration
		// search nodes without parents (head)
		for _, n := range r.NodeOrder.Values() {
			if len(r.Deps[n]) == 0 {
				head = n
				break
			}
		}
		if head == nil {
			// invalid parents specified that are not in nodes
			// add the all rest nodes and exit
			ret = append(ret, r.NodeOrder.Values()...)
			break
		}
		// remove head from all node deps
		for n, _ := range r.Nodes {
			r.Deps[n] = RemoveIf(r.Deps[n], func(e Alteration) bool { return e == head })
		}
		// remove heads from node list
		delete(r.Nodes, head)
		r.NodeOrder.Remove(head)

		// append the nodes
		ret = append(ret, head)
	}

	return ret
}

type LinkedList struct {
	head *LinkedListNode
}

type LinkedListNode struct {
	value  Alteration
	before *LinkedListNode
	next   *LinkedListNode
}

func NewLinkedList(values []Alteration) LinkedList {
	var next *LinkedListNode
	var cur *LinkedListNode
	for i := len(values) - 1; i >= 0; i-- {
		next = cur
		cur = &LinkedListNode{values[i], nil, cur}
		if next != nil {
			next.before = cur
		}
	}
	return LinkedList{cur}
}

func (r LinkedList) Values() []Alteration {
	cur := r.head
	ret := []Alteration{}
	for cur != nil {
		ret = append(ret, cur.value)
		cur = cur.next
	}
	return ret
}

func (r *LinkedList) Remove(value Alteration) {
	cur := r.head
	for cur != nil {
		if cur.value == value {
			if cur.before != nil {
				cur.before.next = cur.next
			} else {
				r.head = cur.next
				if r.head != nil {
					r.head.before = nil
				}
			}
			if cur.next != nil {
				cur.next.before = cur.before
			}
		}
		cur = cur.next
	}
}
