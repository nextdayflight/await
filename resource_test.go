package main

import "testing"

func TestParseResourcesSuccess(t *testing.T) {
	ress := []string{
		"http://user:pass@localhost:1234/path?query=val#fragment",
		"https://user:pass@localhost:1234/path?query=val#fragment",

		"ws://user:pass@localhost:42/path?query=val#fragment",
		"wss://user:pass@localhost:42/path?query=val#fragment",

		"tcp://localhost:42",
		"tcp4://localhost:42",
		"tcp6://[::1]:42",

		"file://relative/path/to/file",
		"file://relative/path/to/file#absent",
		"file:///absolute/path/to/file",
		"file://absolute/path/to/file#absent",

		"postgres://user:pass@localhost:5432/dbname?query=val#fragment",

		"mysql://user:pass@localhost:3306/dbname?query=val#fragment",

		"amqp://host:5672/vhost",
		"amqps://host:5672/vhost",

		"command",
		"command with args",
		"relative/path/to/command",
		"relative/path/to/command with args",
		"/absolute/path/to/command",
		"/absolute/path/to/command with args",

		"",
	}

	actualRess, err := parseResources(ress)
	if err != nil {
		t.Errorf("failed to parse ressources: %#v", err.Error())
	}
	if len(actualRess) != len(ress) {
		t.Errorf("missing parsed ressources")
	}
}

func TestParseResourcesFailure(t *testing.T) {
	_, err := parseResources([]string{"//foo test"})
	if err == nil {
		t.Errorf("expected error parsing invalid ressource")
	}
}

func TestParseFragment(t *testing.T) {
	tests := map[string]map[string]string{
		"":              map[string]string{},
		"=":             map[string]string{}, // invalid format, should be skipped
		"=ignore":       map[string]string{}, // invalid format, should be skipped
		"foo":           map[string]string{"foo": ""},
		"foo=":          map[string]string{"foo": ""},
		"foo=bar":       map[string]string{"foo": "bar"},
		"foo=bar&baz":   map[string]string{"foo": "bar", "baz": ""},
		"foo=bar&baz=1": map[string]string{"foo": "bar", "baz": "1"},
	}
	for given, expected := range tests {
		actual := parseFragment(given)
		if len(actual) != len(expected) {
			t.Errorf("unexpected parsed fragment count %d != %d: %#v",
				len(expected), len(actual), actual)
		}
		for actualK, actualV := range actual {
			if expectedV, ok := expected[actualK]; !ok || actualV[0] != expectedV {
				t.Errorf("unexpected parsed fragment k/v pair")
			}
		}
	}
}
