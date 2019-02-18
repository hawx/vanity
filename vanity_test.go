package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"hawx.me/code/assert"
)

func TestVanity(t *testing.T) {
	server := httptest.NewServer(Server("example.org", Config{
		"/example":     {"/example", "git", "git://example.org/example"},
		"/":            {"/", "git", "https://code.org/r/p/exproj"},
		"/example/pkg": {"/example/pkg", "git", "git://example.org/expkg"},
	}))
	defer server.Close()

	table := map[string]struct {
		req, resp string
	}{
		"simple match": {"/example?go-get=1", "example.org/example git git://example.org/example"},
		"match root":   {"/pkg/foo?go-get=1", "example.org/ git https://code.org/r/p/exproj"},
		"match subpkg": {"/example/pkg/foo?go-get=1", "example.org/example/pkg git git://example.org/expkg"},
	}

	for name, tc := range table {
		t.Run(name, func(t *testing.T) {
			resp, err := http.Get(server.URL + tc.req)
			assert.Nil(t, err)

			body, _ := ioutil.ReadAll(resp.Body)
			assert.Equal(t, `<meta name="go-import" content="`+tc.resp+`">`, string(body))
		})
	}
}

func TestVanityNotGoTool(t *testing.T) {
	server := httptest.NewServer(Server("example.com", Config{
		"/example": {"/example", "git", "git://example.org/example"},
	}))
	defer server.Close()

	resp, err := http.Get(server.URL + "/example")
	assert.Nil(t, err)

	assert.Equal(t, "https://godoc.org/example.com/example", resp.Request.URL.String())
}

func TestReadConfig(t *testing.T) {
	r := strings.NewReader(`/example git git://example.org/example
/what  hg    https://example.com/what`)

	conf, err := DecodeConfig(r)
	if err != nil {
		t.Fatalf("DecodeConfig: %v", err)
	}

	expected := Config{
		"/example": {"/example", "git", "git://example.org/example"},
		"/what":    {"/what", "hg", "https://example.com/what"},
	}

	assert.Equal(t, expected, conf)
}
