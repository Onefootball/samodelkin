// Package gorpmock provides basic data structures
// for mocking gorp related db query logic
package gorpmock

import (
	"github.com/motain/gorp"
)

// DBMapMock is a struct type
// which implements a gorp.SqlExecutor
// interface
// The struct is supposed to be used
// as a test helper for mocking db quering
// logic
type DBMapMock struct {
	gorp.SqlExecutor
	SelectOneFunc func(holder interface{}, query string, args []interface{}) error
	SelectFunc    func(holder interface{}, query string, args []interface{}) ([]interface{}, error)
}

// SelectOne calls m.SelectOneFunc
func (m *DBMapMock) SelectOne(holder interface{}, query string, args ...interface{}) error {
	return m.SelectOneFunc(holder, query, args)
}

// Select calls m.SelectFunc
func (m *DBMapMock) Select(holder interface{}, query string, args ...interface{}) ([]interface{}, error) {
	return m.SelectFunc(holder, query, args)
}
