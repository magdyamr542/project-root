package cmd

// Removing a path by its suffix

type RmCmd struct {
	Suffix string `arg:"" name:"path suffix" help:"Suffix in paths to be deleted"`
}

func (rmCmd *RmCmd) Run() error {
	println("TODO: implement remove path", rmCmd.Suffix)
	return nil
}
