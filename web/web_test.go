package web

import "testing"

func TestEmbeddedFilesExist(t *testing.T) {
	_, err := WebFiles.ReadFile("service-worker.js")
	if err != nil {
		t.Fatalf("failed to read embedded file: %v", err)
	}
}
