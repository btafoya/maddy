/*
Maddy Mail Server - Composable all-in-one email server.
Copyright © 2019-2020 Max Mazurov <fox.cpp@disroot.org>, Maddy Mail Server contributors

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package apprise

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/emersion/go-message/textproto"
	"github.com/emersion/go-smtp"
	"github.com/foxcpp/maddy/framework/buffer"
	"github.com/foxcpp/maddy/framework/config"
	"github.com/foxcpp/maddy/framework/module"
	"github.com/foxcpp/maddy/internal/testutils"
)

func TestAppriseTarget(t *testing.T) {
	t.Parallel()

	// Mock Apprise HTTP server
	mockAppriseServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/" {
			t.Errorf("Expected URL path '/', got %s", r.URL.Path)
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Failed to read request body: %v", err)
		}
		defer r.Body.Close()

		receivedForm, err := url.ParseQuery(string(bodyBytes))
		if err != nil {
			t.Errorf("Failed to parse request body as form: %v", err)
		}

		expectedMessage := "New email from: sender@example.com\nTo: recipient@example.com\nSubject: Test Subject\n\nBody Snippet:\nTest Body"
		if receivedForm.Get("body") != expectedMessage {
			t.Errorf("Received message body mismatch:\nExpected: %s\nGot: %s", expectedMessage, receivedForm.Get("body"))
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer mockAppriseServer.Close()

	// Create Apprise target module
	appriseMod, err := New(modName, "test_apprise", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	m := appriseMod.(*Target)
	m.log = testutils.Logger(t, modName) // Use test logger

	// Configure Apprise URL to point to the mock server
	cfg := config.NewMap(nil, config.Node{
		Children: []config.Node{
			{Name: "apprise_url", Args: []string{mockAppriseServer.URL}},
		},
	})
	if err := m.Init(cfg); err != nil {
		t.Fatal(err)
	}

	// Simulate email delivery
	ctx := context.Background()
	msgMeta := &module.MsgMetadata{}
	mailFrom := "sender@example.com"
	rcptTo := "recipient@example.com"

	delivery, err := m.Start(ctx, msgMeta, mailFrom)
	if err != nil {
		t.Fatal(err)
	}
	defer delivery.Commit(ctx)

	if err := delivery.AddRcpt(ctx, rcptTo, smtp.RcptOptions{}); err != nil {
		t.Fatal(err)
	}

	header := textproto.Header{}
	header.Set("From", mailFrom)
	header.Set("To", rcptTo)
	header.Set("Subject", "Test Subject")

	body := buffer.MemoryBuffer{
		Slice: []byte("Test Body"),
	}

	if err := delivery.Body(ctx, header, body); err != nil {
		t.Fatal(err)
	}
}

func TestAppriseTarget_MultipleURLs(t *testing.T) {
	t.Parallel()

	// Mock Apprise HTTP server for multiple URLs
	callCount := 0
	mockAppriseServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	}))
	defer mockAppriseServer.Close()

	// Create Apprise target module with multiple URLs
	appriseMod, err := New(modName, "test_apprise_multi", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	m := appriseMod.(*Target)
	m.log = testutils.Logger(t, modName)

	cfg := config.NewMap(nil, config.Node{
		Children: []config.Node{
			{Name: "apprise_url", Args: []string{mockAppriseServer.URL, mockAppriseServer.URL + "/second"}},
		},
	})
	if err := m.Init(cfg); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	msgMeta := &module.MsgMetadata{}
	mailFrom := "sender@example.com"
	rcptTo := "recipient@example.com"

	delivery, err := m.Start(ctx, msgMeta, mailFrom)
	if err != nil {
		t.Fatal(err)
	}
	defer delivery.Commit(ctx)

	if err := delivery.AddRcpt(ctx, rcptTo, smtp.RcptOptions{}); err != nil {
		t.Fatal(err)
	}

	header := textproto.Header{}
	header.Set("From", mailFrom)
	header.Set("To", rcptTo)
	header.Set("Subject", "Test Subject")

	body := buffer.MemoryBuffer{
		Slice: []byte("Test Body"),
	}

	if err := delivery.Body(ctx, header, body); err != nil {
		t.Fatal(err)
	}

	if callCount != 2 {
		t.Errorf("Expected Apprise server to be called 2 times, got %d", callCount)
	}
}
