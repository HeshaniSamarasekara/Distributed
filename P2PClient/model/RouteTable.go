package model

// RouteTable - Route table to store joined nodes
type RouteTable struct {
	Nodes []Node
}

// Node - Node struct
type Node struct {
	Name   string
	Port   string
	IP     string
	Status bool
}
