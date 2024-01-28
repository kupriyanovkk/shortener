package encrypt

import (
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	encrypt, err := Get()
	if err != nil {
		t.Errorf("Get() returned an error: %v", err)
	}
	if reflect.DeepEqual(encrypt.Nonce, make([]byte, 12)) {
		t.Errorf("Get() returned a zero nonce")
	}
	if encrypt.AEAD == nil {
		t.Errorf("Get() returned a nil AEAD")
	}
}
