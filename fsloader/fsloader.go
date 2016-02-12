// Package fsloader provides an abstraction
// for reading files kept on the filesystem,
// desiarizing and loading the value into specified
// data types.
//
// Package is meant to be generally used as a config loader
// and follows Onefootball CI/CD constrains:
//  - env var OFCONFIGPATH defines the root dir the package will
//  be scanning for a provided file name
//
//  - in case env var OFCONFIGPATH is not defined: current working directory
//  is taken as the default one
package fsloader

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
)

// OFCONFIGPATH const holds the name
// of the env var which hold the path
// to the application config location
const OFCONFIGPATH = "OFCONFIGPATH"

type (
	// FileLoader is an interface type
	// describing fs file load logic
	FileLoader interface {
		Load(name string, v interface{}) error
	}

	// JSONConfigLoader is a struct type
	// which is used for loading json config files
	JSONConfigLoader struct {
		FileReader
	}

	// ByteConfigLoader is a struct type
	// which is used for loading byte config files
	ByteConfigLoader struct {
		FileReader
	}

	// FileReader is an interface data type
	// which describes fs file read logic
	FileReader interface {
		Read(name string) ([]byte, error)
	}

	// FileReaderFunc is func type and is a
	// FileReader interface adapter
	FileReaderFunc func(name string) ([]byte, error)
)

// Read calls f(name)
func (f FileReaderFunc) Read(name string) ([]byte, error) {
	return f(name)
}

// NewJSONConfigLoader inits and retuns a new
// JSONConfigLoader pointer
func NewJSONConfigLoader() *JSONConfigLoader {
	return &JSONConfigLoader{FileReader: FileReaderFunc(Read)}
}

// Load reads json config file by the name and loads its
// contents into the "v" value
func (cl *JSONConfigLoader) Load(name string, v interface{}) error {
	b, err := cl.Read(name)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, v); err != nil {
		return err
	}

	return nil
}

// NewByteConfigLoader inits and retuns a new ByteConfigLoader
// pointer
func NewByteConfigLoader() *ByteConfigLoader {
	return &ByteConfigLoader{FileReader: FileReaderFunc(Read)}
}

// Load reads string config file and loads its contents
// into "v" value
func (bl *ByteConfigLoader) Load(name string, v interface{}) error {
	b, err := bl.Read(name)
	if err != nil {
		return err
	}

	vb, ok := v.(*[]byte)
	if !ok {
		return errors.New("expected type: *[]byte")
	}

	// write bytes into the "v" value
	*vb = b
	return nil
}

// Read reads file contents from the fs
// the file directory is either a project root or
// is read from env var defined under OFCONFIGPATH
func Read(name string) ([]byte, error) {
	filekey := path.Join(getDir(), name)
	return ioutil.ReadFile(filekey)
}

// getDir either a dir defined under env var OFCONFIGPATH
// or project root path
func getDir() string {
	if dir := os.Getenv(OFCONFIGPATH); dir != "" {
		return dir
	}

	dir, _ := os.Getwd()
	return dir
}
