package main

import (
	"project-root/src/cmd"
	"project-root/src/fs"

	"github.com/alecthomas/kong"
)

var CLI struct {
	Add   cmd.AddCmd   `cmd:"" help:"Register the given path as the root of the project (can be relative)"`
	Ls    cmd.LsCmd    `cmd:"" help:"List all saved project roots" aliases:"l,ls,list"`
	Clear cmd.ClearCmd `cmd:"" help:"clears the database of saved projects. Will delete everything. Use with CAUTION" aliases:"clr"`
	Go    cmd.GoCmd    `cmd:"" help:"go to the root of this project. (can be omitted which means 'pr=pr go')" default:"1"`
	Purge cmd.PurgeCmd `cmd:"" help:"delete all registered paths that no longer exist in the file system"`
	Back  cmd.BackCmd  `cmd:"" help:"go back to the place you went to the project root from" aliases:"b"`
	To    cmd.ToCmd    `cmd:"" help:"go to a specific project by using a saved path prefix" aliases:"t"`
}

func main() {
	ctx := kong.Parse(&CLI)
	ctx.BindTo(fs.DefaultFileSystem, (*fs.FileSystemHandler)(nil))
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
