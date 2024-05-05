package inspect

import (
	"os"
	"path/filepath"

	"github.com/imkouga/kratos-pkg/inspect/alarmer"
)

const (
	separate = "|#|"
)

func Init(metric string) error {

	path, err := os.Executable()
	_, exec := filepath.Split(path)
	if nil != err {
		return err
	}

	return Init1(exec, metric)
}

func Init1(servicName, metric string) error {

	if len(metric) <= 0 {
		return ErrInitParamMustBeProvided
	}

	if err := alarmer.Init(metric); nil != err {
		return err
	}

	setServiceName(servicName)
	setOnOff(true)

	initBasic()
	return nil
}
