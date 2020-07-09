package git

import (
	"context"
	"os/exec"
	"time"

	"github.com/Oppodelldog/pulli/log"
)

const defaultExecutionTimeout = 20 * time.Second

type execWithTimeoutFuncDef func(ctx context.Context, s1 string, s2 ...string) *exec.Cmd

var execFunc = execWithTimeoutFuncDef(exec.CommandContext)

func git(dir string, s ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultExecutionTimeout)
	defer cancel()

	defer func() {
		err := ctx.Err()
		if err != nil {
			log.Printf("error executing git: %v", err)
		}
	}()

	cmd := execFunc(ctx, "git", s...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()

	return string(output), err
}
