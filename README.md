goRailsYourself
===============

A suite of random functions needed when porting/mixing Go/Rails code. Use at your own risk, don't expect much and feel free to send lots of awesome pull requests.


## Usage:

    inflector.Parameterize("Matt Aïmonetti", "-") => "matt-aimonetti"


See the [documentation](http://godoc.org/github.com/mattetti/goRailsYourself) and/or the test suite for more examples.

## Dependencies:

The inflector package relies on
[unidecode](http://godoc.org/github.com/fiam/gounidecode/unidecode) to
handle the transliteration.

The test suite uses
[Goblin](http://tech.gilt.com/post/64409561192/goblin-a-minimal-and-beautiful-testing-framework-for)


