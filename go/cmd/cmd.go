package cmd

// Adding a new path
type AddCmd struct {
	Path string `arg:"" name:"path" help:"Path to add" type:"path"`
}

func (addCmd *AddCmd) Run() error {
	return RegisterProject(addCmd.Path)
}

// Listing all paths
type LsCmd struct{}

func (lsCmd *LsCmd) Run() error {
	println("TODO: implement ls paths")
	return nil
}

// Removing a path by its suffix

type RmCmd struct {
	Suffix string `arg:"" name:"path suffix" help:"Suffix in paths to be deleted"`
}

func (rmCmd *RmCmd) Run() error {
	println("TODO: implement remove path", rmCmd.Suffix)
	return nil
}

// Clearing the db

type ClearCmd struct{}

func (clrCmd *ClearCmd) Run() error {
	println("TODO: implement clear db")
	return nil
}

// Going to the root

type GoCmd struct{}

func (goCmd *GoCmd) Run() error {
	println("TODO: implement go")
	return nil
}

// Purge
type PurgeCmd struct{}

func (purgeCmd *PurgeCmd) Run() error {
	println("TODO: implement purge")
	return nil
}
