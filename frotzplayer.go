package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type frotzPlayer struct {
	cmd       *exec.Cmd
	bufferOut *bytes.Buffer
	cmdOut    io.Reader
	// cmdErr    io.ReadCloser
	cmdIn io.WriteCloser
	story string
}

func (fp frotzPlayer) IsRunning() bool {
	if fp.cmd == nil {
		return false
	}
	if fp.cmd.ProcessState != nil && !fp.cmd.ProcessState.Exited() {
		return true
	}
	return fp.cmd.Process != nil && fp.cmd.Process.Pid > 0
}

func (fp *frotzPlayer) Run() error {
	return fp.cmd.Run()
}

func (fp *frotzPlayer) Input(in string) error {
	in = strings.TrimSpace(strings.Split(in, "\n")[0]) + "\n"
	_, err := fp.cmdIn.Write([]byte(in))
	return err
}

func (fp *frotzPlayer) Output() (string, error, error) {
	buf := make([]byte, 1024)
	if fp.cmdOut == nil {
		return "", nil, nil
	}
	var out string
	_, outErr := fp.cmdOut.Read(buf)
	out = string(buf)
	fmt.Println(out)

	return out, outErr, nil /* errErr */
}

func (fp *frotzPlayer) Kill() error {
	if fp.cmd != nil && fp.cmd.Process != nil {
		return fp.cmd.Process.Kill()
	}
	return nil
}

func newFrotzPlayer(frotzPath string, gamePath string) (*frotzPlayer, error) {
	_, err := os.Stat(gamePath)
	if err != nil {
		return nil, err
	}

	player := new(frotzPlayer)
	player.bufferOut = new(bytes.Buffer)
	player.cmd = exec.Command(frotzPath, gamePath)

	stdOut, err := player.cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stdErr, err := player.cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	player.cmdOut = io.MultiReader(stdOut, stdErr)

	if player.cmdIn, err = player.cmd.StdinPipe(); err != nil {
		return nil, err
	}

	return player, err
}
