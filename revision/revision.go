// Package revision provides a convenience abstraction
// for reading repository file with REVISION information.
//
// REVISION file pattern is a contrained designed for Onefootball CI/CD
// environment in order to keep track of the current app version.
// CD (continious delivery) implies project build and deploy on merge event.
// REVISION file is expected to be generated within the Makefile build process
// and deployed to the production server.
//
// It is a developer's responsibility to decide about the contents of the REVISION file,
// as well as how to make file contents accessible.
//
// Onefootball Use Case
//
// Sample incorporation flow
//
//  - echo latest git commit hash to REVISION (Makefile): @echo `git log -n 1 --pretty=format:"%H"` > REVISION
//  - add the file to deploy package: tar -jcf ./app.bz2 *
//  - make REVISION file contents available from the /_healthcheck_ endpoint
package revision

import (
	"github.com/onefootball/samodelkin/fsloader"
)

// revisionFile holds file name
// which stores the latest git revision id
const revisionFile = "REVISION"

// AppRevision is a []byte type
// for loading and storing REVISION file contents
type AppRevision []byte

// Message returns a byte slice
// which holds the application revision id
// message and can be used to pass to io.Writer
func (r AppRevision) Message() []byte {
	return append([]byte("revision: "), r...)
}

// String returns AppRevision string
// representation
func (r AppRevision) String() string {
	return string(r)
}

// Load reads REVISION file contents
// and stores the value to the method receiver
func (r *AppRevision) Load() error {
	bl := fsloader.NewByteConfigLoader()

	var b []byte
	if err := bl.Load(revisionFile, &b); err != nil {
		return err
	}

	*r = AppRevision(b)
	return nil
}
