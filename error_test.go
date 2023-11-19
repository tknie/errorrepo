/*
* Copyright 2023 Thorsten A. Knieling
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
 */

package errorrepo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tknie/log"
)

const logPrefix = "LOG:"

type OutputLog struct {
}

func (*OutputLog) Debugf(format string, args ...interface{}) {
	fmt.Printf(logPrefix+format+"\n", args...)
}

func (*OutputLog) Infof(format string, args ...interface{}) {
	fmt.Printf(logPrefix+format+"\n", args...)
}

func (*OutputLog) Errorf(format string, args ...interface{}) {
	fmt.Printf(logPrefix+format+"\n", args...)
}

func (*OutputLog) Fatalf(format string, args ...interface{}) {
	fmt.Printf(logPrefix+format+"\n", args...)
}

func (*OutputLog) Fatal(args ...interface{}) {
	fmt.Printf(logPrefix+"%#v", args)
}

func TestMessages(t *testing.T) {
	log.Log = &OutputLog{}
	msg := &Message{text: "abc"}
	assert.Equal(t, "abc", msg.convertArgs())
	msg.text = "bcd {0} {1}"
	assert.Equal(t, "bcd a 1", msg.convertArgs("a", 1))
	msg.text = "cde %v %d"
	assert.Equal(t, "cde b 1", msg.convertArgs("b", 1))
}

func TestErrors(t *testing.T) {
	log.Log = &OutputLog{}
	err := NewError("ERR00001")
	assert.Nil(t, err.(*Error).err)
	assert.Equal(t, "ERR00001: ID not found", err.Error())
	err = NewError("ERR00004", "xyz")
	assert.Nil(t, err.(*Error).err)
	assert.Equal(t, "ERR00004: Test message xyz", err.Error())
	err = fmt.Errorf("abc")
	xerr := NewError("ERR00003", err)
	assert.Equal(t, "ERR00003: Error db open: abc", xerr.Error())
	assert.Error(t, xerr.(*Error).err)
	assert.Equal(t, err, xerr.(*Error).err)
	err = NewError("ERR65535")
	assert.Equal(t, "ERR65535: not implemented", err.Error())
	err = NewError("ERR50001", "testing")
	assert.Equal(t, "ERR50001: Internal error: testing", err.Error())
	err = NewError("ERR50002", "testing2")
	assert.Equal(t, "ERR50002: Internal error parameter: testing2", err.Error())
}

func TestRegisterErrors(t *testing.T) {
	RegisterMessage("en",
		`ABC00001=2323232
DEV001111=abc dfdjfdkj {0} {1}
`)
	err := NewError("ABC00001")
	assert.Nil(t, err.(*Error).err)
	assert.Equal(t, "ABC00001: 2323232", err.Error())

	err = NewError("DEV001111", "par1", "par2", "par3")
	assert.Nil(t, err.(*Error).err)
	assert.Equal(t, "DEV001111: abc dfdjfdkj par1 par2%!(EXTRA string=par3)", err.Error())

}
