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

	uuid "github.com/satori/go.uuid"
)

type ConsoleManager struct {
	Consoles map[uuid.UUID]*Console
}

func (cm *ConsoleManager) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			for _, con := range cm.Consoles {
				if !con.Closed() {
					con.Close()
				}
			}

			break
		default:
		}

		for _, con := range cm.Consoles {
			if con.Closed() {
				delete(cm.Consoles, con.UUID)
			}
		}
	}
}

func (cm *ConsoleManager) Add(con *Console) {
	cm.Consoles[con.UUID] = con
}

func (cm *ConsoleManager) Get(uid uuid.UUID) (*Console, bool) {
	con, ok := cm.Consoles[uid]
	return con, ok
}

func (cm *ConsoleManager) NewConsole(loader Loader) (*Console, error) {
	con, err := StartConsole(loader)
	if err != nil {
		return nil, err
	}

	cm.Add(con)

	return con, nil
}

func StartConsole(loader Loader) (*Console, error) {
	con := &Console{}

	var err error

	con.UUID, err = uuid.NewV4()
	if err != nil {
		return nil, err
	}

	con.Cmder = &Cmder{
		WorkingDir: loader.Path(),
	}

	program, args := loader.Cmd()

	con.Cmder.Start(program, args...)

	go con.start()

	return con, nil
}

type Console struct { // sync for each session
	UUID  uuid.UUID
	Cmder *Cmder

	Logs *LogStacker

	closeCh chan bool
	closed  bool
}

func (con *Console) Closed() bool {
	return con.closed
}

func (con *Console) start() {
	for {
		line, ok := con.Cmder.Line()
		if !ok {
			break
		}

		con.Logs.Add(line)
	}
}

func (con *Console) Close() {
	if con.closed {
		return
	}

	con.closed = true

	con.Cmder.Close()
}

func (con *Console) Lines(t *LogTracker) []string {
	return con.Logs.AllChanges(t)
}

func (con *Console) SendCommand(cmd string) {
	con.Cmder.Send(cmd)
}

func NewLogTracker() *LogTracker {
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

func (st *LogStacker) Changes(n int, t *LogTracker) []string {
	t.ChangeCounter = 0

	if t.ChangeCounter == -1 {
		return st.Get(MinInt(n, MinInt(st.counter, st.logs.Len())))
	}

	return st.Get(MinInt(n, MinInt(t.ChangeCounter, st.logs.Len())))
}

func (st *LogStacker) AllChanges(t *LogTracker) []string {
	t.ChangeCounter = 0

	return st.Get(MinInt(st.counter, st.logs.Len()))
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

func (st *LogStacker) add(str string) {
	st.logs.Value = str
	st.logs = st.logs.Next()
}

func (st *LogStacker) prev(n int) []string {
	logs := st.logs // it won't change for master

	var ok bool
	data := make([]string, n)
	for i := 0; i < n; i++ {
		logs = logs.Prev()
		data[i], ok = logs.Value.(string)
		if !ok {
			panic("couldn't convert to string")
		}
	}

	return data
}
