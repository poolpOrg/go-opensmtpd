/*
 * Copyright (c) 2020 Gilles Chehade <gilles@poolp.org>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package table

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type LookupService string

const (
	ALIAS       LookupService = "alias"       /* returns struct expand	*/
	DOMAIN      LookupService = "domain"      /* returns struct destination	*/
	CREDENTIALS LookupService = "credentials" /* returns struct credentials	*/
	NETADDR     LookupService = "netaddr"     /* returns struct netaddr	*/
	USERINFO    LookupService = "userinfo"    /* returns struct userinfo	*/
	SOURCE      LookupService = "source"      /* returns struct source	*/
	MAILADDR    LookupService = "mailaddr"    /* returns struct mailaddr	*/
	ADDRNAME    LookupService = "addrname"    /* returns struct addrname	*/
	MAILADDRMAP LookupService = "mailaddrmap" /* returns struct maddrmap	*/
	RELAYHOST   LookupService = "relayhost"   /* returns struct relayhost	*/
	STRING      LookupService = "string"
	REGEX       LookupService = "regex"
)

func nameToLookupService(name string) LookupService {
	switch name {
	case "alias":
		return ALIAS
	case "domain":
		return DOMAIN
	case "credentials":
		return CREDENTIALS
	case "netaddr":
		return NETADDR
	case "userinfo":
		return USERINFO
	case "source":
		return SOURCE
	case "mailaddr":
		return MAILADDR
	case "addrname":
		return ADDRNAME
	case "mailaddrmap":
		return MAILADDRMAP
	case "relayhost":
		return RELAYHOST
	case "string":
		return STRING
	case "regex":
		return REGEX
	}
	return STRING
}

type updatePrototype func(string, string)
type checkPrototype func(string, string, LookupService, string)
type lookupPrototype func(string, string, LookupService, string)
type fetchPrototype func(string, string, LookupService)

var updateCb updatePrototype = Updated
var checkCb checkPrototype
var lookupCb lookupPrototype
var fetchCb fetchPrototype

func Failure(token string) {
	fmt.Printf("table-result|%s|failure\n", token)
}

func Updated(token string) {
	fmt.Printf("table-result|%s|updated\n", token)
}

func Boolean(token string, result bool) {
	if result == false {
		fmt.Printf("table-result|%s|not-found\n", token)
	} else {
		fmt.Printf("table-result|%s|found\n", token)
	}
}

func Result(token string, result ...string) {
	if len(result) == 0 {
		fmt.Printf("table-result|%s|not-found\n", token)
	} else {
		fmt.Printf("table-result|%s|found|%s\n", token, result[0])
	}
}

func OnUpdate(cb updatePrototype) {
	updateCb = cb
}

func OnCheck(cb checkPrototype) {
	checkCb = cb
}

func OnLookup(cb lookupPrototype) {
	lookupCb = cb
}

func OnFetch(cb fetchPrototype) {
	fetchCb = cb
}

func Dispatch() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			os.Exit(0)
		}

		line := scanner.Text()
		atoms := strings.Split(line, "|")
		if len(atoms) < 6 {
			log.Fatalf("missing atoms: %s", line)
		}

		if atoms[0] != "table" {
			log.Fatalf("invalid stream: %s", atoms[0])
		}

		if atoms[1] != "0.1" {
			log.Fatalf("unsupported protocol version: %s", atoms[1])
		}

		table := atoms[3]
		token := atoms[5]
		switch atoms[4] {
		case "update":
			updateCb(token)

		case "check":
			checkCb(token, nameToLookupService(atoms[6]), strings.Join(atoms[7:], "|"))

		case "lookup":
			lookupCb(token, nameToLookupService(atoms[6]), strings.Join(atoms[7:], "|"))

		case "fetch":
			fetchCb(token, nameToLookupService(atoms[6]))

		default:
			log.Fatalf("unsupported operation: %s", atoms[4])
		}
	}
}
