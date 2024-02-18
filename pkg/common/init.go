// pkg/common/init.go
package common

import (
	"encoding/gob"
)

func init() {
	gob.Register(Order{})
	gob.Register(User{})
	// Register other types if needed
}
