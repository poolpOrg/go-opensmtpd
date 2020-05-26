package main

import (
	"./table"
)

func update(token string) {
	table.Updated(token)
	//	table.Failure(token)
}

func check(token string, service table.LookupService, key string) {
	table.Boolean(token, true)
	//table.Boolean(token, false)
	//table.Failure(token)
}

func lookup(token string, service table.LookupService, key string) {
	table.Result(token, "foobar")
	//table.Failure(token)
}

func fetch(token string, service table.LookupService) {
	table.Result(token, "foobar")
	//table.Failure(token)
}

func main() {
	table.OnUpdate(update)
	table.OnCheck(check)
	table.OnLookup(lookup)
	table.OnFetch(fetch)
	table.Dispatch()
}
