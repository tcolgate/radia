package graphalg

import "testing"

func TestRegister1(t *testing.T) {
	RegisterMessage(TestMessage{})
}

func TestRegister2(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Register of non proto.Message did not panic")
		}
	}()
	type other struct{}
	RegisterMessage(other{})
}
