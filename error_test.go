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

package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {

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
	assert.Equal(t, "DEV001111: abc dfdjfdkj par1 par2", err.Error())

}
