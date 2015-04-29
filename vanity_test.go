package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestVanity(t *testing.T) {
	server := httptest.NewServer(Server("example.org", Config{
		{"/example", "git", "git://example.org/example"},
		{"/", "git", "https://code.org/r/p/exproj"},
		{"/example/pkg", "git", "git://example.org/expkg"},
	}))
	defer server.Close()

	table := []struct {
		req, resp string
	}{
		{"/example?go-get=1", "example.org/example git git://example.org/example"},
		{"/pkg/foo?go-get=1", "example.org/ git https://code.org/r/p/exproj"},
		{"/example/pkg/foo?go-get=1", "example.org/example/pkg git git://example.org/expkg"},
	}

	for _, tc := range table {
		resp, err := http.Get(server.URL + tc.req)
		if err != nil {
			t.Fatalf("Get: %v", err)
		}

		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) != `<meta name="go-import" content="`+tc.resp+`">` {
			t.Fatalf("body incorrect, got: %s", body)
		}
	}
}

func TestVanityNotGoTool(t *testing.T) {
	server := httptest.NewServer(Server("example.com", Config{
		{"/example", "git", "git://example.org/example"},
	}))
	defer server.Close()

	resp, err := http.Get(server.URL + "/example")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if resp.Request.URL.String() != "http://godoc.org/example.com/example" {
		t.Fatalf("Expected redirect, Got: %s", resp.Request.URL)
	}
}

func TestReadConfig(t *testing.T) {
	r := strings.NewReader(`/example git git://example.org/example
/what  hg    https://example.com/what`)

	conf, err := DecodeConfig(r)
	if err != nil {
		t.Fatalf("DecodeConfig: %v", err)
	}

	expected := Config{
		{"/example", "git", "git://example.org/example"},
		{"/what", "hg", "https://example.com/what"},
	}

	if len(conf) != 2 || conf[0] != expected[0] || conf[1] != expected[1] {
		t.Fatalf("Expected: %v, Got: %v", expected, conf)
	}
}
