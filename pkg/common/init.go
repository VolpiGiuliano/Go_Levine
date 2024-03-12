// pkg/common/init.go
package common

import (
	"encoding/gob"
)

// Put all the variables that are sent
func init() {
	gob.Register(Order{})
	gob.Register(User{})
	gob.Register(Order_Book{})
	gob.Register(Queue{})
	
}
