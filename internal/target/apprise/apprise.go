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
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/emersion/go-message/textproto"
	"github.com/emersion/go-smtp"
	"github.com/foxcpp/maddy/framework/buffer"
	"github.com/foxcpp/maddy/framework/config"
	"github.com/foxcpp/maddy/framework/log"
	"github.com/foxcpp/maddy/framework/module"
	"github.com/foxcpp/maddy/internal/target"
)

const modName = "target.apprise"

type Target struct {
	instName string
	appriseURLs []string // Stores configured Apprise URLs
	log         log.Logger
}

func New(_, instName string, _, inlineArgs []string) (module.Module, error) {
	t := &Target{
		instName: instName,
		log:      log.Logger{Name: modName},
	}
	t.appriseURLs = inlineArgs // Apprise URLs can be passed as inline arguments

	return t, nil
}

func (t *Target) Init(cfg *config.Map) error {
	var urlsFromCfg []string
	cfg.Bool("debug", true, false, &t.log.Debug)
	cfg.StringList("apprise_url", false, false, nil, &urlsFromCfg)

	if _, err := cfg.Process(); err != nil {
		return err
	}

	if len(t.appriseURLs) == 0 {
		t.appriseURLs = urlsFromCfg
	} else if len(urlsFromCfg) > 0 {
		t.appriseURLs = append(t.appriseURLs, urlsFromCfg...)
	}

	if len(t.appriseURLs) == 0 {
		return fmt.Errorf("%s: at least one Apprise URL is required", modName)
	}

	// Validate URLs
	for _, appriseURL := range t.appriseURLs {
		if _, err := url.ParseRequestURI(appriseURL); err != nil {
			return fmt.Errorf("%s: invalid Apprise URL '%s': %w", modName, appriseURL, err)
		}
	}

	return nil
}

func (t *Target) Name() string {
	return modName
}

func (t *Target) InstanceName() string {
	return t.instName
}

func (t *Target) Close() error {
	return nil
}

type delivery struct {
	t        *Target
	msgMeta  *module.MsgMetadata
	mailFrom string
	rcptTo   []string // Store recipients for notification
	log      log.Logger
}

func (t *Target) Start(ctx context.Context, msgMeta *module.MsgMetadata, mailFrom string) (module.Delivery, error) {
	return &delivery{
		t:        t,
		msgMeta:  msgMeta,
		mailFrom: mailFrom,
		log:      target.DeliveryLogger(t.log, msgMeta),
	}, nil
}

func (d *delivery) AddRcpt(ctx context.Context, rcptTo string, opts smtp.RcptOptions) error {
	d.rcptTo = append(d.rcptTo, rcptTo)
	return nil
}

func (d *delivery) Body(ctx context.Context, header textproto.Header, body buffer.Buffer) error {
	subject := header.Get("Subject")
	from := header.Get("From")
	to := header.Get("To") // Corrected: use header.Get("To")

	// Construct a basic notification message
	message := fmt.Sprintf("New email from: %s\nTo: %s\nSubject: %s", from, to, subject)
	
	// Optionally read a snippet of the body
	var bodySnippet string
	if r, err := body.Open(); err == nil {
		buf := new(bytes.Buffer)
		_, err := io.CopyN(buf, r, 512) // Read up to 512 bytes
		r.Close()
		if err != nil && err != io.EOF {
			d.log.Println("Failed to read body snippet:", err) // Corrected logging
		}
		bodySnippet = buf.String()
		if len(bodySnippet) > 0 {
			message += "\n\nBody Snippet:\n" + bodySnippet
		}
	} else {
		d.log.Println("Failed to open message body:", err) // Corrected logging
	}


	for _, appriseURL := range d.t.appriseURLs {
		form := url.Values{}
		form.Add("body", message)
		
		req, err := http.NewRequestWithContext(ctx, "POST", appriseURL, strings.NewReader(form.Encode()))
		if err != nil {
			d.log.Error("Failed to create Apprise request", err, "url", appriseURL)
			continue
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			d.log.Error("Failed to send Apprise notification", err, "url", appriseURL)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			respBody, _ := io.ReadAll(resp.Body)
			d.log.Error("Apprise notification failed with non-2xx status", nil, "url", appriseURL, "status", resp.Status, "response_body", string(respBody))
		} else {
			if d.t.log.Debug { // Guard with debug boolean
				d.log.Debugf("Apprise notification sent successfully (URL: %s)", appriseURL) // Corrected logging
			}
		}
	}

	return nil // Do not block email delivery if notification fails
}

func (d *delivery) Abort(ctx context.Context) error {
	return nil
}

func (d *delivery) Commit(ctx context.Context) error {
	return nil
}

func init() {
	module.Register(modName, New)
}
