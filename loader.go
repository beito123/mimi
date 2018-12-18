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
	"errors"
	"path/filepath"
)

type LoaderManager struct {
	Loaders map[string]Loader
}

func (lm *LoaderManager) Get(name string) (Loader, bool) {
	loader, ok := lm.Loaders[name]
	if !ok {
		return nil, false
	}

	return loader.New(), true
}

func (lm *LoaderManager) Add(name string, loader Loader) {
	lm.Loaders[name] = loader
}

func (lm *LoaderManager) Remove(name string) {
	delete(lm.Loaders, name)
}

type Loader interface {
	Path() string
	Cmd() (string, []string)

	Init(path string) error

	New() Loader
}

type PMMPLoader struct {
	path string
}

func (loader *PMMPLoader) Path() string {
	return loader.path
}

func (loader *PMMPLoader) Init(path string) error {
	loader.path = filepath.Clean(path)

	// check
	if ExistFile(loader.Program()) {
		return errors.New("Couldn't find php program")
	}

	if ExistFile(loader.Target()) {
		return errors.New("Couldn't find PMMP program")
	}

	return nil
}

func (loader *PMMPLoader) Program() string {
	if IsWin() {
		return loader.path + "/bin/php/php.exe"
	}

	return loader.path + "/bin/php/php"
}

func (loader *PMMPLoader) Target() string {
	if ExistFile(loader.path + "/src/pocketmine/PocketMine.php") { // TODO: support to change
		return loader.path + "/src/pocketmine/PocketMine.php"
	}

	return loader.path + "/PocketMine-MP.phar"
}

func (loader *PMMPLoader) Cmd() (string, []string) {
	return loader.Program(), []string{"-c", "bin/php", loader.Target()}
}

func (PMMPLoader) New() Loader {
	return new(PMMPLoader)
}
