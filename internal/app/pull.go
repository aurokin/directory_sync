package app

import "io"

func runPull(args []string, stdout, stderr io.Writer) int {
	return runSync("pull", args, stdout, stderr)
}
