package main

import (
	"project-root/src/cmd"

	"github.com/alecthomas/kong"
)

var CLI struct {
	Add   cmd.AddCmd   `cmd:"" help:"Register the given path as the root of the project (can be relative)"`
	Ls    cmd.LsCmd    `cmd:"" help:"List all saved project roots" aliases:"l,ls"`
	Clear cmd.ClearCmd `cmd:"" help:"clears the database of saved projects. Will delete everything. Use with CAUTION" aliases:"clr"`
	Go    cmd.GoCmd    `cmd:"" help:"go to the root of this project" default:"1"`
	Purge cmd.PurgeCmd `cmd:"" help:"delete all registered paths that no longer exist in the file system"`
}

func main() {
	ctx := kong.Parse(&CLI)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
