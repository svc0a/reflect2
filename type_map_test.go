package reflect2

import "testing"

type User struct {
}

func TestRegister(t *testing.T) {
	Register[User]()
}
