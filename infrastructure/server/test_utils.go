package server

import "testing"

func assertResponseStatus(t *testing.T, want, got int) {
	t.Helper()
	if got != want {
		t.Errorf("\nincorrect http status\nwant: %d, got %d", want, got)
	}
}
