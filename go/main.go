package main

import (
	"fmt"
	"project-root/fs"
)

func main() {
	dir, err := fs.GetHomeDir()
	fmt.Println(dir, err)
}
