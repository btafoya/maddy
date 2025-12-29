//go:build !cgo || !libpam
// +build !cgo !libpam

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

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "maddy-pam-helper: PAM support not available")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "This binary was built without PAM support. To enable PAM authentication:")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "1. Install PAM development headers:")
	fmt.Fprintln(os.Stderr, "   - Debian/Ubuntu: apt install libpam0g-dev")
	fmt.Fprintln(os.Stderr, "   - Fedora/RHEL:   dnf install pam-devel")
	fmt.Fprintln(os.Stderr, "   - Alpine:        apk add linux-pam-dev")
	fmt.Fprintln(os.Stderr, "   - Arch:          pacman -S pam")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "2. Rebuild with the libpam tag:")
	fmt.Fprintln(os.Stderr, "   go build -tags libpam ./cmd/maddy-pam-helper")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "See: https://maddy.email/reference/auth/pam/")
	os.Exit(2)
}
