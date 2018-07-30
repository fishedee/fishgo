package command

func Help(argv []string) (string, error) {
	return `
FishCmd is a tool for golang build, it is awesome fast!!!

Usage:

	fishcmd command [arguments]

The commands are:
	
	clean			Clean a go application
	run	[appName]	Run a go application
		--watch		AutoRun a go application when dictory file change
	test 			Test a go application
		--watch		AutoTest a go application when dictory file change
		--benchmark	Benchmark a go application when dictory file change
	version			FishCmd version
	help			FishCmd help

`, nil
}
