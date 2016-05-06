package command

func Help(argv []string) (string, error) {
	return `
FishCmd is a tool for golang build, it is awesome fast!!!

Usage:

	fishcmd command [arguments]

The commands are:
	
	clean		Clean a go application
	build		Build a go application
		-fast 	Build fast more and more
	watch		AutoBuild a go application when dictory file change
		-fast 	Build fast more and more
	version		FishCmd version
	help		FishCmd help

`, nil
}
