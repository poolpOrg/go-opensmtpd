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

type updatePrototype func(string)
type checkPrototype func(string, LookupService, string)
type lookupPrototype func(string, LookupService, string)
type fetchPrototype func(string, LookupService)

var updateCb updatePrototype
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
		if len(atoms) < 5 {
			log.Fatalf("missing atoms: %s", line)
		}

		if atoms[0] != "table" {
			log.Fatalf("invalid stream: %s", atoms[0])
		}

		if atoms[1] != "0.1" {
			log.Fatalf("unsupported protocol version: %s", atoms[1])
		}

		token := atoms[4]
		switch atoms[3] {
		case "update":
			updateCb(token)

		case "check":
			checkCb(token, nameToLookupService(atoms[5]), strings.Join(atoms[6:], "|"))

		case "lookup":
			lookupCb(token, nameToLookupService(atoms[5]), strings.Join(atoms[6:], "|"))

		case "fetch":
			fetchCb(token, nameToLookupService(atoms[5]))

		default:
			log.Fatalf("unsupported operation: %s", atoms[3])
		}
	}
}
