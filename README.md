#Parse

[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/f0ster/parse) [![license](http://img.shields.io/badge/license-BSD-red.svg?style=flat)](https://raw.githubusercontent.com/f0ster/parse/master/LICENSE)


This is a go lib for using parse-server, it has been unofficially tested with the android and iOS parse SDKs, and parse for nodejs.  

###Installation

    go get github.com/f0ster/parse

###Usage:
```go
package main

import (
    "fmt"
	"time"
    
    "github.com/f0ster/parse"
)

func main() {
    // the last argument is a mount point where parse-server is used an express app, at the specified host
    parse.Initialize(viper.GetString("parse_app_id"), viper.GetString("parse_rest_key"), viper.GetString("parse_master_key"), "parse-server-hostname.coolwebsite.com", "http", "/parse")
    parse.SetUserAgent("myapi.coolwebsite.com")
    

    user := parse.User{}
    q, err := parse.NewQuery(&user)
	if err != nil {
		panic(err)
	}
    q.EqualTo("email", "cooluser@gmail.com")
    q.GreaterThan("numFollowers", 10).OrderBy("-createdAt") // API is chainable
    err := q.First()
    if err != nil {
        if pe, ok := err.(parse.ParseError); ok {
            fmt.Printf("Error querying parse: %d - %s\n", pe.Code(), pe.Message())
        }
    }
    
    fmt.Printf("Retrieved user with id: %s\n", u.Id)

	q2, _ := parse.NewQuery(&parse.User{})
	q2.GreaterThan("createdAt", time.Date(2014, 01, 01, 0, 0, 0, 0, time.UTC))

	rc := make(chan *parse.User)

	// .Each will retrieve all results for a query and send them to the provided channel
	// The iterator returned allows for early cancelation of the iteration process, and
	// stores any error that triggers early termination
	iterator, err := q2.Each(rc)
	for u := range rc{
		fmt.Printf("received user: %v\n", u)
		// Do something
		if err := process(u); err != nil {
			// Cancel if there was an error
			iterator.Cancel()
		}
	}

	// An error occurred - not all rows were processed
	if it.Error() != nil {
		panic(it.Error())
	}
}
```

###TODO
- Add structured logging
- Missing query operations
	- Related to
- Missing CRUD operations:
    - Update
		- Field ops (__op):
			- AddRelation
			- RemoveRelation
- Roles
- Cloud Functions
- Background Jobs
- Analytics
- File upload/retrieval
- Batch operations
