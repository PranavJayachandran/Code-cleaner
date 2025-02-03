package domain

import (
	"sort"
)

type Node struct {
	Node       NodeType
	Value      string
	FileName   string
	Dependents []string
	Callers    []string
  LineNumber int
  IsExported bool
}

func (Node *Node) CreateNode(nodeType NodeType, value string, fileName string, dependents []string, lineNumber int, isExported bool) {
	Node.Node = nodeType
	Node.Value = value
	Node.FileName = fileName
	Node.Dependents = dependents
	Node.Callers = []string{}
  Node.LineNumber = lineNumber
  Node.IsExported = isExported
}

func (Node Node) Equals(node Node) bool {
	if node.Node != Node.Node {
		return false
	}
	nodeDependencies := node.Dependents
	i, j := 0, 0
	for i < len(Node.Dependents) && j < len(nodeDependencies) {
		if Node.Dependents[i] != nodeDependencies[j] {
			return false
		}
	}
	return true
}

func (Node *Node) AddDenpendents(dependency string) {
	Node.Dependents = append(Node.Dependents, dependency)
	sort.Slice(Node.Dependents, func(i, j int) bool {
		return Node.Dependents[i] < Node.Dependents[j]
	})
}
