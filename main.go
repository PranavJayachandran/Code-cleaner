package main

import (
	"code_cleaner/pkg/file"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
  startDir, _ := os.Getwd()
  if len(os.Args) > 1 {
    startDir = os.Args[1]
  }
  fmt.Printf("Starting the walk from %s\n", startDir);
  file.InitialiseWalks(startDir)
  startDir = filepath.Join(startDir, "src")
  err := filepath.Walk(startDir, file.DeclarationWalk);
  if err != nil{
    fmt.Println(err)
  }
  err = filepath.Walk(startDir, file.ImportWalk);
  if err != nil {
    fmt.Println(err)
  }
  file.Result()
}
