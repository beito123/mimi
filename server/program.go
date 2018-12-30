package server

/*
 * mimi
 *
 * Copyright (c) 2018 beito
 *
 * This software is released under the MIT License.
 * http://opensource.org/licenses/mit-license.php
**/

import (
	"errors"
	"regexp"
	"strings"

	"gitlab.com/beito123/mimi/config"
)

var RegOnlyAlphabetNumber = regexp.MustCompile("^[a-zA-Z0-9]+$")

func NewProgramManager(lm *LoaderManager, programs []config.ProgramConfig) (*ProgramManager, error) {
	pm := &ProgramManager{
		Programs: make(map[string]*Program),
	}

	for _, pc := range programs {
		name := strings.ToLower(pc.Name)

		// Vaild name of a program
		if !RegOnlyAlphabetNumber.Match([]byte(name)) {
			return nil, errors.New("can't use for program's name except alphabets and numbers")
		}

		loader, ok := lm.Get(pc.Loader)
		if !ok {
			return nil, errors.New("couldn't find a loader")
		}

		err := loader.Init(pc.Path, pc.LoaderOptions)
		if err != nil {
			return nil, err
		}

		pm.Add(&Program{
			Name:   name,
			Loader: loader,
		})
	}

	return pm, nil
}

type ProgramManager struct {
	Programs map[string]*Program
}

func (pm *ProgramManager) Get(name string) (*Program, bool) {
	p, ok := pm.Programs[name]

	return p, ok
}

func (pm *ProgramManager) Add(p *Program) {
	pm.Programs[p.Name] = p
}

func (pm *ProgramManager) Remove(name string) {
	delete(pm.Programs, name)
}

type Program struct {
	Name   string
	Loader Loader
}
