package kmgCmd

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type KmgCmd struct {
	*exec.Cmd
}

func FromString(cmdS string) *KmgCmd {
	return &KmgCmd{execFromString(cmdS)}
}

func (cmd *KmgCmd) SetDir(dir string) *KmgCmd {
	cmd.Dir = dir
	return cmd
}

type MuplitCmd struct {
	CmdList []*exec.Cmd
}

func NewMuplitCmd(cmds []string) *MuplitCmd {
	mc := &MuplitCmd{}
	for _, v := range cmds {
		mc.CmdList = append(mc.CmdList, execFromString(v))
	}
	return mc
}
func (cmds *MuplitCmd) SetDir(dir string) *MuplitCmd {
	for _, cmd := range cmds.CmdList {
		cmd.Dir = dir
	}
	return cmds
}
func (cmds *MuplitCmd) SetStdOut(w io.Writer) *MuplitCmd {
	for _, cmd := range cmds.CmdList {
		cmd.Stdout = w
	}
	return cmds
}
func (cmds *MuplitCmd) SetStdErr(w io.Writer) *MuplitCmd {
	for _, cmd := range cmds.CmdList {
		cmd.Stderr = w
	}
	return cmds
}
func (cmds *MuplitCmd) SerializeRun() (err error) {
	for i, cmd := range cmds.CmdList {
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("SerializeRun fail, index:%d err:%s", i, err.Error())
		}
	}
	return nil
}

func execFromString(cmdS string) *exec.Cmd {
	args := strings.Split(cmdS, " ")
	return exec.Command(args[0], args[1:]...)
}
