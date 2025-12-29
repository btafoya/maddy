package dkim

import (
	"context"
	"testing"

	"github.com/emersion/go-message/textproto"
	"github.com/foxcpp/maddy/framework/buffer"
	"github.com/foxcpp/maddy/framework/config"
	"github.com/foxcpp/maddy/framework/module"
	"github.com/foxcpp/maddy/internal/testutils"
)

func TestDKIM_Table(t *testing.T) {
	dkimMod, err := New("modify.dkim", "test", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	m := dkimMod.(*Modifier)
	m.keyPathTemplate = "../../../tests/testdata/dkim_keys/{domain}_{selector}.key"
	m.selector = "default"

	// Mock the domainsTable
	m.domainsTable = testutils.Table{M: map[string]string{
		"example.org": "",
		"foobar.com":  "",
	}}

	// Initialize the modifier (it will skip loading keys if domainsTable is set)
	if err := m.Init(config.NewMap(nil, config.Node{})); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	test := func(domain string, shouldSign bool) {
		t.Helper()
		state, err := m.ModStateForMsg(ctx, &module.MsgMetadata{})
		if err != nil {
			t.Fatal(err)
		}
		defer state.Close()

		if _, err := state.RewriteSender(ctx, "test@"+domain); err != nil {
			t.Fatal(err)
		}

		hdr := textproto.Header{}
		hdr.Add("From", "test@"+domain)
		body := buffer.MemoryBuffer{
			Slice: []byte("test"),
		}

		if err := state.RewriteBody(ctx, &hdr, &body); err != nil {
			t.Fatal(err)
		}

		sig := hdr.Get("DKIM-Signature")
		if shouldSign {
			if sig == "" {
				t.Error("expected signature, got none")
			}
		} else {
			if sig != "" {
				t.Error("expected no signature, got one:", sig)
			}
		}
	}

	test("example.org", true)
	test("foobar.com", true)
	test("google.com", false)
}
