package app

import "github.com/linn221/memory-sheets/models"

// TheSession is a global session for the application
// since there is only one user
var TheSession models.Session

func init() {
	TheSession = models.Session{}
}
