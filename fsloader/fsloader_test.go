package fsloader

import (
	"encoding/json"
	"errors"
	"testing"

	gc "github.com/motain/gocheck"
)

type FSLoaderTestSuite struct{}

var _ = gc.Suite(&FSLoaderTestSuite{})

func TestFSLoader(t *testing.T) { gc.TestingT(t) }

func (s *FSLoaderTestSuite) TestJSONConfigLoaderSuccess(c *gc.C) {
	expectedJson := `{"a":"b"}`
	expectedName := "testname"

	fr := FileReaderFunc(func(name string) ([]byte, error) {
		c.Check(name, gc.Equals, expectedName)
		return []byte(expectedJson), nil
	})

	jsonLoader := &JSONConfigLoader{FileReader: fr}

	var v interface{}
	err := jsonLoader.Load(expectedName, &v)
	c.Check(err, gc.IsNil)

	obtainedJson, err := json.Marshal(v)
	c.Check(err, gc.IsNil)
	c.Check(string(obtainedJson), gc.Equals, expectedJson)
}

func (s *FSLoaderTestSuite) TestJSONConfigLoaderReadError(c *gc.C) {
	expectedErr := errors.New("test error")

	fr := FileReaderFunc(func(name string) ([]byte, error) {
		return nil, expectedErr
	})

	var v interface{}

	jsonLoader := &JSONConfigLoader{FileReader: fr}
	c.Check(jsonLoader.Load("test", &v), gc.Equals, expectedErr)
}

func (s *FSLoaderTestSuite) TestJSONConfigLoaderUnmarshalError(c *gc.C) {
	expectedJson := `{"a":"b"}`
	expectedName := "testname"

	fr := FileReaderFunc(func(name string) ([]byte, error) {
		c.Check(name, gc.Equals, expectedName)
		return []byte(expectedJson), nil
	})

	jsonLoader := &JSONConfigLoader{FileReader: fr}

	var v interface{}
	err := jsonLoader.Load(expectedName, v)
	c.Check(err, gc.NotNil)
}

func (s *FSLoaderTestSuite) TestByteConfigLoaderSuccess(c *gc.C) {
	expectedBody := []byte("test config")
	expectedName := "testname"

	fr := FileReaderFunc(func(name string) ([]byte, error) {
		c.Check(name, gc.Equals, expectedName)
		return expectedBody, nil
	})

	byteLoader := &ByteConfigLoader{FileReader: fr}

	var v []byte
	err := byteLoader.Load(expectedName, &v)
	c.Check(err, gc.IsNil)
	c.Check(string(v), gc.Equals, string(expectedBody))
}

func (s *FSLoaderTestSuite) TestByteConfigLoaderReadError(c *gc.C) {
	expectedErr := errors.New("test error")

	fr := FileReaderFunc(func(name string) ([]byte, error) {
		return nil, expectedErr
	})

	byteLoader := &ByteConfigLoader{FileReader: fr}

	var v []byte
	err := byteLoader.Load("test", &v)
	c.Check(err, gc.Equals, expectedErr)
}

func (s *FSLoaderTestSuite) TestByteConfigLoaderTypeError(c *gc.C) {
	fr := FileReaderFunc(func(name string) ([]byte, error) {
		return []byte("test config"), nil
	})

	byteLoader := &ByteConfigLoader{FileReader: fr}

	var v interface{}
	err := byteLoader.Load("test", &v)
	c.Check(err.Error(), gc.Equals, "expected type: *[]byte")
}
