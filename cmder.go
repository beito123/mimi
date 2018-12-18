package mimi

/*
 * mimi
 *
 * Copyright (c) 2018 beito
 *
 * This software is released under the MIT License.
 * http://opensource.org/licenses/mit-license.php
**/

import (
	"bufio"
	"io"
	"os/exec"
)

type Cmder struct {
	WorkingDir string

	lineCh chan string
	errCh  chan error
	sendCh chan string
}

func (cmder *Cmder) Start(program string, param ...string) error {
	cmder.lineCh = make(chan string, 10)
	cmder.sendCh = make(chan string, 10)
	cmder.errCh = make(chan error, 1)

	cmd := exec.Command(program, param...)

	cmd.Dir = cmder.WorkingDir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	/* TODO:
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	*/

	err = cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		defer cmder.close()

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			cmder.lineCh <- scanner.Text()
		}

		err = cmd.Wait()
		if err != nil {
			cmder.errCh <- err
		}
	}()

	go func() {
		select {
		case str, ok := <-cmder.sendCh:
			if ok {
				return
			}

			io.WriteString(stdin, str)
		}
	}()

	return nil
}

func (cmder *Cmder) Close() {
	cmder.close()
}

func (cmder *Cmder) close() {
	close(cmder.errCh)
	close(cmder.lineCh)
	close(cmder.sendCh)
}

func (cmder *Cmder) Line() (string, bool) {
	select {
	case err, ok := <-cmder.errCh:
		if ok {
			return err.Error(), true // ummm...
		}
	case line, ok := <-cmder.lineCh:
		return line, ok
	}

	return "", false
}

func (cmder *Cmder) Send(str string) {
	cmder.sendCh <- str
}
