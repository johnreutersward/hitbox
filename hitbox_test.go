package hitbox

import (
	"testing"
)

func TestGames(t *testing.T) {
	hbc := NewClient(nil)
	_, _, err := hbc.Games()
	if err != nil {
		t.Errorf("expected no error, got = %s", err.Error())
	}
}
