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
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/tknie/log"
)

const defaultLanguage = "en"

type Message struct {
	id   string
	text string
}

var localesMap map[string]map[string]*Message

type Error struct {
	id  string
	msg string
	err error
}

//go:embed messages
var embedFiles embed.FS

// init create message default message hash
func init() {
	localesMap = make(map[string]map[string]*Message)
	fss, err := embedFiles.ReadDir("messages")
	if err != nil {
		panic("Internal config load error: " + err.Error())
	}
	RegisterDirectory(fss)
}

func RegisterDirectory(fss []fs.DirEntry) {
	for _, f := range fss {
		if f.Type().IsRegular() {
			byteValue, err := embedFiles.ReadFile("messages/" + f.Name())
			if err != nil {
				panic("Internal config load error: " + err.Error())
			}
			lang := path.Ext(f.Name())
			RegisterMessage(lang[1:], string(byteValue))
		}
	}
}

// language current message language
func language() string {
	lang := os.Getenv("LANG")
	switch {
	case lang == "":
		lang = "en"
	default:
		if len(lang) < 2 {
			lang = "en"
		} else {
			lang = lang[0:2]
		}
	}
	log.Log.Debugf("Current LANG: %s", lang)
	return lang
}

// RegisterMessage register message for locale to be registered
// in the current error list
func RegisterMessage(locale, msgData string) error {
	lines := strings.Split(msgData, "\n")
	var messageHash map[string]*Message
	var ok bool
	if messageHash, ok = localesMap[locale]; !ok {
		messageHash = make(map[string]*Message)
		localesMap[locale] = messageHash
	}
	for _, line := range lines {
		index := strings.IndexByte(line, '=')
		if index == -1 {
			return fmt.Errorf("message structure parse error")
		}
		id := line[:index]
		text := line[index+1:]
		log.Log.Debugf("Register %s -> %s", locale, id)
		messageHash[id] = &Message{id, text}
	}
	return nil
}

// NewErrorLocale create new message with given parameter
// locale locale language given
// errID error ID of current error
// args argument for the error
func NewErrorLocale(locale, errID string, args ...interface{}) error {
	var err error
	if len(args) > 0 {
		switch e := args[len(args)-1].(type) {
		case error:
			err = e
		default:
		}
	}
	e := &Error{id: errID, err: err}
	e.createMessage(locale, args...)
	log.Log.Debugf("Error %s created: %v", errID, err)
	log.Log.Debugf("Error message created:[%s] %s", errID, e.msg)
	log.Log.Debugf("Stack trace:\n%s", string(debug.Stack()))

	return e
}

// NewError create new error with system language
func NewError(errID string, args ...interface{}) error {
	return NewErrorLocale(language(), errID, args...)
}

// ID provide error ID
func (e *Error) ID() string {
	return e.id
}

// Message provide the error message only
func (e *Error) Message() string {
	return e.msg
}

// Error provide error interface
func (e *Error) Error() string {
	return fmt.Sprintf("%8s: %s", e.id, e.msg)
}

// createMessage create a new message with locale and given args
func (e *Error) createMessage(locale string, args ...interface{}) {
	messageHash := localesMap[locale]
	if messageHash == nil {
		messageHash = localesMap[language()]
	}
	if messageHash == nil {
		messageHash = localesMap[defaultLanguage]
	}
	var outLine *Message
	if cmsg, ok := messageHash[e.id]; ok {
		outLine = cmsg
	}
	log.Log.Debugf("Get %s: %s", e.id, outLine)
	if outLine != nil {
		log.Log.Debugf("Search %s: %s", e.id, outLine)
		m := outLine.text
		if len(args) > 0 {
			m = outLine.convertArgs(args...)
		}
		// e.msg = fmt.Sprintf("%8s: %s", e.id, m)
		e.msg = m
	} else {
		e.msg = fmt.Sprintf("Unknown error ...%s", e.id)
	}
}

// convertArgs define replace arguments
func (m *Message) convertArgs(args ...interface{}) string {
	msg := m.text
	c := `\{\d\}`
	re := regexp.MustCompile(c)
	msg = re.ReplaceAllString(msg, "%v")
	return fmt.Sprintf(msg, args...)
}
