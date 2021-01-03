goRailsYourself
===============

[![GoDoc](http://godoc.org/github.com/mattetti/goRailsYourself?status.png)](https://pkg.go.dev/github.com/mattetti/goRailsYourself)


A suite of packages useful when you have to deal with Go and Rails apps
or when migrating from Ruby to Go. 

The crypto package allows for shared authentication cookie support with Rails, included version 5.2+.


See the [documentation](http://godoc.org/github.com/mattetti/goRailsYourself) and/or the test suite for more examples.

## Dependencies:

The inflector package relies on:
 [unidecode](http://godoc.org/github.com/fiam/gounidecode/unidecode) to handle the transliteration.

The crypto package relies on:
  [pbkdf2](http://golang.org/x/crypto/pbkdf2) to handle the
generation of derived keys.

The test suite uses
[Goblin](http://tech.gilt.com/post/64409561192/goblin-a-minimal-and-beautiful-testing-framework-for)


