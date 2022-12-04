package server

import "testing"

func TestViewersId(t *testing.T) {
	path := "/posts/concurrency-abstractions-in-go/"
	expected := "views_posts/concurrency-abstractions-in-go.md"
	actual, ok := viewersId(path)
	if !ok {
		t.Errorf("viewersId(%s) returned false", path)
	}
	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}
