package mimi

import (
	"regexp"
	"runtime"
)

// Error

type ErrorReceiver interface {
	HandleError(error)
}

type ErrorHandler struct {
	receivers []ErrorReceiver
}

func (hand *ErrorHandler) Handle(err error) {
	for _, recv := range hand.receivers {
		recv.HandleError(err)
	}
}

// Dumper

const MaxStacksCount = 30

func Dump(err error) {
	stacks := Stack(1, MaxStacksCount)

	logger.Fatalf("Happened errors: %s\n\n" + err.Error())

	count := len(stacks)
	for i, s := range stacks {
		logger.Fatalf("%d: %s@%s:%s:%d\n", count-i, s.PackageName, s.FileName, s.FunctionName, s.FileLine)
	}
}

//Thank you, sgykfjsm from http://sgykfjsm.github.io/blog/2016/01/20/golang-function-tracing/
var regStack = regexp.MustCompile(`^(\S.+)\.(\S.+)$`)

type CallerInfo struct {
	PackageName  string
	FunctionName string
	FileName     string
	FileLine     int
}

func Stack(skip int, count int) (callerInfo []*CallerInfo) {
	for i := 1; i <= count; i++ {
		pc, _, _, ok := runtime.Caller(skip + i)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc)
		fileName, fileLine := fn.FileLine(pc)

		_fn := regStack.FindStringSubmatch(fn.Name())
		callerInfo = append(callerInfo, &CallerInfo{
			PackageName:  _fn[1],
			FunctionName: _fn[2],
			FileName:     fileName,
			FileLine:     fileLine,
		})
	}

	return callerInfo
}
