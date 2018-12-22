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
	"strings"
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

func (lm *LoaderManager) Add(loader Loader) {
	lm.Loaders[loader.Name()] = loader
}

func (lm *LoaderManager) Remove(name string) {
	delete(lm.Loaders, name)
}

type Loader interface {
	Name() string
	Path() string
	Cmd() (string, []string)

	Init(path string, options map[string]string) error

	New() Loader
}

type PMMPLoader struct {
	path     string
	PHPPath  string
	MainPath string
	Args     []string
}

func (PMMPLoader) Name() string {
	return "PMMP"
}

func (loader *PMMPLoader) Path() string {
	return loader.path
}

func (loader *PMMPLoader) Init(path string, options map[string]string) (err error) {
	loader.path, err = filepath.Abs(filepath.Clean(path))
	if err != nil {
		return err
	}

	phpPath, ok := options["php"]
	if ok {
		loader.PHPPath = phpPath
	} else {
		if IsWin() {
			loader.PHPPath = loader.path + "/bin/php/php.exe"
		} else {
			loader.PHPPath = loader.path + "/bin/php/php"
		}
	}

	mainPath, ok := options["main"]
	if ok {
		loader.MainPath = mainPath
	} else {
		if ExistFile(loader.path + "/src/pocketmine/PocketMine.php") {
			loader.MainPath = loader.path + "/src/pocketmine/PocketMine.php"
		} else {
			loader.MainPath = loader.path + "/PocketMine-MP.phar"
		}
	}

	args, ok := options["args"]
	if ok {
		loader.Args = strings.Fields(args)
	} else {
		loader.Args = []string{"-c", "bin/php"}
	}

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
	return loader.PHPPath
}

func (loader *PMMPLoader) Target() string {
	return loader.MainPath
}

func (loader *PMMPLoader) Cmd() (string, []string) {
	return loader.Program(), append(loader.Args, loader.Target())
}

func (PMMPLoader) New() Loader {
	return new(PMMPLoader)
}
