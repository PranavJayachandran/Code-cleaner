package domain

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

func NodeFinder(line string, filePath string, lineNumber int) (Node, error) {
	fmt.Println(filePath)
	if isTrue, val, isExport := isInterface(line); isTrue {
		node := Node{}
		ext := filepath.Ext(filePath)
		filePath = strings.TrimSuffix(filePath, ext)
		node.CreateNode(Interface, val, filePath, []string{}, lineNumber, isExport)
		return node, nil
	}
	return Node{}, errors.New("No matching node for this type")
}

func isInterface(line string) (bool, string, bool) {
	elements := strings.Split(line, " ")
	if elements[0] == "export" && elements[1] == "interface" {
		return true, strings.TrimRight(strings.TrimSpace(elements[2]), "{"), true
	}
	return false, "", false
}
