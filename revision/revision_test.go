package revision

import (
	"testing"

	gc "github.com/go-check/check"
)

type RevisionTestSuite struct{}

var _ = gc.Suite(&RevisionTestSuite{})

func TestRevision(t *testing.T) { gc.TestingT(t) }

func (s *RevisionTestSuite) TestAppRevisionString(c *gc.C) {
	r := AppRevision([]byte("123456"))
	c.Check(r.String(), gc.Equals, "123456")
}

func (s *RevisionTestSuite) TestAppRevisionMessage(c *gc.C) {
	r := AppRevision([]byte("123456"))
	c.Check(string(r.Message()), gc.Equals, "revision: 123456")
}
