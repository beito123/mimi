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
	"container/ring"
	"context"
)

type ConsoleManager struct {
	Consoles map[string]*Console
}

func StartConsole(loader Loader) *Console {
	con := &Console{}

	con.Cmder = &Cmder{
		WorkingDir: loader.Path(),
	}

	program, args := loader.Cmd()

	con.Cmder.Start(program, args...)

	return con
}

type Console struct { // sync for each session
	Cmder *Cmder

	Logs *LogStacker

	closeCh chan bool
	closed  bool
}

func (con *Console) start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			break
		}

		con.Cmder.Line()
	}
}

func (con *Console) Close() {
	if con.closed {
		return
	}

	con.closed = true

	con.Cmder.Close()
}

func (con *Console) Line() (string, bool) {
	return con.Cmder.Line() // change
}

func (con *Console) SendCommand(cmd string) {
	con.Cmder.Send(cmd)
}

func NewLogTracker(n int) *LogTracker {
	return &LogTracker{
		ChangeCounter: -1,
	}
}

type LogTracker struct {
	id            int
	ChangeCounter int
}

func NewLogStacker(n int) *LogStacker {
	return &LogStacker{
		logs: ring.New(n),
	}
}

type LogStacker struct {
	logs    *ring.Ring
	counter int

	trackerID int
	trackers  []*LogTracker
}

func (st *LogStacker) AddTracker(t *LogTracker) {
	t.id = st.trackerID
	st.trackerID++

	st.trackers = append(st.trackers, t)
}

func (st *LogStacker) RemoveTracker(t *LogTracker) {
	nt := make([]*LogTracker, len(st.trackers))
	for _, v := range st.trackers {
		if v.id != t.id {
			nt = append(nt, v)
		}
	}
}

func (st *LogStacker) Add(str string) {
	st.add(str)

	st.counter++

	for _, t := range st.trackers {
		t.ChangeCounter++
	}
}

func (st *LogStacker) Get(n int) []string {
	return st.prev(n)
}

func (st *LogStacker) GetChanges(n int, t *LogTracker) []string {
	if t.ChangeCounter == -1 {
		return st.Get(MinInt(n, MinInt(st.counter, st.logs.Len())))
	}

	return st.Get(MinInt(n, MinInt(t.ChangeCounter, st.logs.Len())))
}

func (st *LogStacker) add(str string) {
	st.logs.Value = str
	st.logs = st.logs.Next()
}

func (st *LogStacker) prev(n int) []string {
	logs := st.logs // it won't change for master

	data := make([]string, n)
	for i := 0; i < n; i++ {
		logs = logs.Prev()
		data[i] = logs.Value.(string)
	}

	return data
}
