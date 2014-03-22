package kmgCmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func SetCmdEnv(cmd *exec.Cmd, key string, value string) error {
	if len(cmd.Env) == 0 {
		cmd.Env = os.Environ()
	}
	env, err := NewEnvFromArray(cmd.Env)
	if err != nil {
		return err
	}
	env.Values[key] = value
	cmd.Env = env.ToArray()
	return nil
}

type Env struct {
	Values map[string]string
}

func NewEnvFromArray(env []string) (envObj *Env, err error) {
	envObj = &Env{Values: make(map[string]string)}
	for _, v1 := range env {
		pos := strings.IndexRune(v1, '=')
		if pos == -1 {
			return nil, fmt.Errorf("NewEnvFromArray: input string not have =, string: %s", v1)
		}
		key := v1[:pos]
		v2 := v1[pos+1:]

		_, ok := envObj.Values[key]
		if ok {
			//ignore this error and use first value
			continue
			//return nil,errors.Sprintf("NewEnvFromArray: two keys has same name, name: %s",key)
		}
		envObj.Values[key] = v2
	}
	return
}

func (env *Env) ToArray() []string {
	output := make([]string, 0, len(env.Values))
	for k, v := range env.Values {
		output = append(output, k+"="+v)
	}
	return output
}
