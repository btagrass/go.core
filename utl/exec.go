package utl

import (
	"fmt"
	"os/exec"
)

// 命令
func Command(nameArgs ...string) (string, error) {
	var output []byte
	var err error
	for _, na := range nameArgs {
		nas := Split(na, ' ')
		output, err = exec.Command(nas[0], nas[1:]...).CombinedOutput()
		if err != nil {
			return "", fmt.Errorf(string(output))
		}
	}

	return string(output), err
}
