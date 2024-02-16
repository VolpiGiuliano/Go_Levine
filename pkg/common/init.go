// pkg/common/init.go
package common

import (
	"encoding/gob"
)

func init() {
	gob.Register(Order{})
	// Register other types if needed
}
