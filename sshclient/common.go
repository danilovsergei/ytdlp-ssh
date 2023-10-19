package sshclient

import (
	"fmt"
	"log"
	"strconv"
)

func (c *Connection) IsFileOrDirExists(dir string) bool {
	retryFunc := func(output string) bool { return output == "" }

	out := c.ExecCommandWithOutput(fmt.Sprintf("[ -d '%s' ] && echo true || echo false", dir), retryFunc)
	result, err := strconv.ParseBool(out)
	if err != nil {
		log.Fatalf("Failed to parse ssh out %s %s", out, err)
	}
	return result
}
