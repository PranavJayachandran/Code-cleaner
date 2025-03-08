package file

import (
	"code_cleaner/pkg/domain"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type importStatement struct {
	name     string
	filePath string
}

var config Config
var StartPath string

func handleImportStatement(data []byte, i int, path string) []importStatement {
	//First it should be a space of {
	for i < len(data) {
		if data[i] != ' ' && data[i] != '{' {
			return []importStatement{}
		}
		if data[i] == '{' {
			break
		}
		i++
	}

	i++

	if i >= len(data) {
		return []importStatement{}
	}
	importedElements := ""
	//Go till the closing }
	for i < len(data) && data[i] != '}' {
		importedElements += string(data[i])
		i++
	}
	filePath := ""
	i++
	if i >= len(data) {
		return []importStatement{}
	}

	//Go till the end of filePath
	insidePath := false
	for i < len(data) {
		if data[i] == '"' || data[i] == '\'' {
			if insidePath {
				break
			}
			insidePath = true
			i++
			continue
		}
		if insidePath {
			filePath += string(data[i])
		}
		i++
	}
	if !insidePath {
		return []importStatement{}
	}
	if filePath[0] == '@' {
		alias := ""
		j := 0
		for j < len(filePath) && filePath[j] != '/' {
			alias += string(filePath[j])
			j++
		}
		alias += "/*"
		if val, exists := config.CompilerOptions.Paths[alias]; exists {
			filePath = StartPath + "/" + val[0][:len(val[0])-2] + filePath[j:]
		}
	} else {
		filePath = filepath.Clean(filepath.Join(path, "../"+filePath))
	}
	importStatements := []importStatement{}
	importTokens := strings.Split(importedElements, ",")
	for _, val := range importTokens {
		val = strings.Trim(val, "")
		filePath = strings.Replace(filePath, "/", "\\", -1)
		importStatements = append(importStatements, importStatement{filePath: filePath, name: removeNextLineAndSpace(val)})
	}
	return importStatements
}
func removeNextLineAndSpace(text string) string {
	if len(text) == 0 {
		return text
	}
	i := 0
	for i < len(text) && (text[i] == '\n' || text[i] == ' ') {
		i++
	}
	j := len(text) - 1
	for j >= 0 && (text[j] == '\n' || text[j] == ';' || text[j] == ' ') {
		j--
	}
	if i >= j {
		return ""
	}
	return text[i : j+1]
}
func ImportWalk(path string, info os.FileInfo, err error) error {
	if !info.IsDir() {
		data, err := os.ReadFile(path)
		if err == nil {
			for i := 0; i < len(data); i++ {
				if len(data)-i > 6 {
					if string(data[i:i+6]) == "import" {
						i += 7
						importStatements := handleImportStatement(data, i, path)
						// fmt.Printf("\n%s\n", path)
						// for _, el := range importStatements {
						// 	fmt.Println(el)
						// }
						for _, val := range importStatements {
							nodeIndex, exists := nodeMapper[val.name+"__"+val.filePath]
							if exists {
								nodeList[nodeIndex].Callers = append(nodeList[nodeIndex].Callers, path)
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func Result() {
	for _, item := range nodeList {
		fmt.Printf("%d %s %s\n", len(item.Callers), item.Value, item.FileName)
	}
}

var nodeList []domain.Node
var nodeMapper map[string]int

func InitialiseWalks(path string) {
	findAliasPaths(path)
	nodeList = []domain.Node{}
	nodeMapper = make(map[string]int)
}
func DeclarationWalk(path string, info os.FileInfo, err error) error {
	if !info.IsDir() {
		data, err := os.ReadFile(path)
		if err == nil {
			line := ""
			lineCount := 0
			for _, item := range data {
				if item == '\n' || item == ';' {
					line = strings.Trim(line, " ")
					node, err := domain.NodeFinder(line, path, lineCount)
					if err == nil {
						nodeMapper[node.Value+"__"+node.FileName] = len(nodeList)
						nodeList = append(nodeList, node)
					}
					line = ""
					lineCount++
				} else {
					line += string(item)
				}
			}
		}
	}
	return nil
}

type Config struct {
	CompilerOptions Options `json:"compilerOptions"`
}

type Options struct {
	Paths map[string][]string `json:"paths"`
}

// Function to remove comments from JSON data
func removeJSONComments(data string) string {
	// Remove single-line comments (// ...)
	singleLineComment := regexp.MustCompile(`(?m)^\s*//.*$`)
	data = singleLineComment.ReplaceAllString(data, "")

	// Remove multi-line comments (/* ... */)
	multiLineComment := regexp.MustCompile(`(?s)/\*.*?\*/`)
	data = multiLineComment.ReplaceAllString(data, "")

	return data
}

func findAliasPaths(path string) {
	path = filepath.Join(path, "tsconfig.json")

	// Read the JSON file
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading tsconfig.json:", err)
		return
	}

	// Convert data to string and remove comments
	cleanData := removeJSONComments(string(data))

	err = json.Unmarshal([]byte(cleanData), &config)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
}
