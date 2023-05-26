package src

import "testing"

const TestToken = "Bearer 499c6deb356c2b5f6223ebb5a07c6e190b1bc911"

func TestEmptyValidator(t *testing.T) {
	if ok, _ := EmptyValidator("test", ""); !ok {
		t.Error("Empty validator should always return ok")
	}

	if ok, _ := testValidator("test", ""); !ok {
		t.Error("Test validator credentials failed (user test)")
	}

	if ok, _ := testValidator("test2", ""); ok {
		t.Error("Test validator credentials failed, should be false (user test2)")
	}

	if ok, _ := censorTrackerValidator("_", TestToken); !ok {
		t.Error("Censor tracker auth failed")
	}
}
