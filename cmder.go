package mimi

import (
	"bufio"
	"os/exec"
)

type Cmder struct {
	Dir    string
	LineCh chan string
	ErrCh  chan error
}

func (cmder *Cmder) Start(program string, param ...string) error {
	cmder.LineCh = make(chan string, 10)
	cmder.ErrCh = make(chan error, 1)

	cmd := exec.Command(program, param...)

	cmd.Dir = cmder.Dir

	stdout, err := cmd.StdoutPipe()
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
			cmder.LineCh <- scanner.Text()
		}

		err = cmd.Wait()
		if err != nil {
			cmder.ErrCh <- err
		}
	}()

	return nil
}

func (cmder *Cmder) Close() {
	cmder.close()
}

func (cmder *Cmder) close() {
	close(cmder.ErrCh)
	close(cmder.LineCh)
}

func (cmder *Cmder) Line() (string, bool) {
	select {
	case err, ok := <-cmder.ErrCh:
		if ok {
			return err.Error(), true // ummm...
		}
	case line, ok := <-cmder.LineCh:
		return line, ok
	}

	return "", false
}
