package app

import "io"

func runPush(args []string, stdout, stderr io.Writer) int {
	return runSync("push", args, stdout, stderr)
}
