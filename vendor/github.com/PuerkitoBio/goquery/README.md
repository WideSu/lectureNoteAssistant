# goquery - a little like that j-thing, only in Go
[![build status](https://secure.travis-ci.org/PuerkitoBio/goquery.svg?branch=master)](http://travis-ci.org/PuerkitoBio/goquery) [![GoDoc](https://godoc.org/github.com/PuerkitoBio/goquery?status.png)](http://godoc.org/github.com/PuerkitoBio/goquery) [![Sourcegraph Badge](https://sourcegraph.com/github.com/PuerkitoBio/goquery/-/badge.svg)](https://sourcegraph.com/github.com/PuerkitoBio/goquery?badge)

goquery brings a syntax and a set of features similar to [jQuery][] to the [Go language][go]. It is based on Go's [net/html package][html] and the CSS Selector library [cascadia][]. Since the net/html parser returns nodes, and not a full-featured DOM tree, jQuery's stateful manipulation functions (like height(), css(), detach()) have been left off.

Also, because the net/html parser requires UTF-8 encoding, so does goquery: it is the caller's responsibility to ensure that the source document provides UTF-8 encoded HTML. See the [wiki][] for various options to do this.

Syntax-wise, it is as close as possible to jQuery, with the same function names when possible, and that warm and fuzzy chainable interface. jQuery being the ultra-popular library that it is, I felt that writing a similar HTML-manipulating library was better to follow its API than to start anew (in the same spirit as Go's `fmt` package), even though some of its methods are less than intuitive (looking at you, [index()][index]...).

## Table of Contents

* [Installation](#installation)
* [Changelog](#changelog)
* [API](#api)
* [Examples](#examples)
* [Related Projects](#related-projects)
* [Support](#support)
* [License](#license)

## Installation

Please note that because of the net/html dependency, goquery requires Go1.1+.

    $ go get github.com/PuerkitoBio/goquery

(optional) To run unit tests:

    $ cd $GOPATH/src/github.com/PuerkitoBio/goquery
    $ go test

(optional) To run benchmarks (warning: it runs for a few minutes):

    $ cd $GOPATH/src/github.com/PuerkitoBio/goquery
    $ go test -bench=".*"

## API

goquery exposes two structs, `Document` and `Selection`, and the `Matcher` interface. Unlike jQuery, which is loaded as part of a DOM document, and thus acts on its containing document, goquery doesn't know which HTML document to act upon. So it needs to be told, and that's what the `Document` type is for. It holds the root document node as the initial Selection value to manipulate.

jQuery often has many variants for the same function (no argument, a selector string argument, a jQuery object argument, a DOM element argument, ...). Instead of exposing the same features in goquery as a single method with variadic empty interface arguments, statically-typed signatures are used following this naming convention:

*    When the jQuery equivalent can be called with no argument, it has the same name as jQuery for the no argument signature (e.g.: `Prev()`), and the version with a selector string argument is called `XxxFiltered()` (e.g.: `PrevFiltered()`)
*    When the jQuery equivalent **requires** one argument, the same name as jQuery is used for the selector string version (e.g.: `Is()`)
*    The signatures accepting a jQuery object as argument are defined in goquery as `XxxSelection()` and take a `*Selection` object as argument (e.g.: `FilterSelection()`)
*    The signatures accepting a DOM element as argument in jQuery are defined in goquery as `XxxNodes()` and take a variadic argument of type `*html.Node` (e.g.: `FilterNodes()`)
*    The signatures accepting a function as argument in jQuery are defined in goquery as `XxxFunction()` and take a function as argument (e.g.: `FilterFunction()`)
*    The goquery methods that can be called with a selector string have a corresponding version that take a `Matcher` interface and are defined as `XxxMatcher()` (e.g.: `IsMatcher()`)

Utility functions that are not in jQuery but are useful in Go are implemented as functions (that take a `*Selection` as parameter), to avoid a potential naming clash on the `*Selection`'s methods (reserved for jQuery-equivalent behaviour).

The complete [godoc reference documentation can be found here][doc].

Please note that Cascadia's selectors do not necessarily match all supported selectors of jQuery (Sizzle). See the [cascadia project][cascadia] for details. Invalid selector strings compile to a `Matcher` that fails to match any node. Behaviour of the various functions that take a selector string as argument follows from that fact, e.g. (where `~` is an invalid selector string):

* `Find("~")` returns an empty selection because the selector string doesn't match anything.
* `Add("~")` returns a new selection that holds the same nodes as the original selection, because it didn't add any node (selector string didn't match anything).
* `ParentsFiltered("~")` returns an empty selection because the selector string doesn't match anything.
* `ParentsUntil("~")` returns all parents of the selection because the selector string didn't match any element to stop before the top element.

## Examples

See some tips and tricks in the [wiki][].

Adapted from example_test.go:

```Go
package main

import (
  "fmt"
  "log"
  "net/http"

  "github.com/PuerkitoBio/goquery"
)

func ExampleScrape() {
  // Request the HTML page.
  res, err := http.Get("http://metalsucks.net")
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()
  if res.StatusCode != 200 {
    log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
  }

  // Load the HTML document
  doc, err := goquery.NewDocumentFromReader(res.Body)
  if err != nil {
    log.Fatal(err)
  }

  // Find the review items
  doc.Find(".sidebar-reviews article .content-block").Each(func(i int, s *goquery.Selection) {
    // For each item found, get the band and title
    band := s.Find("a").Text()
    title := s.Find("i").Text()
    fmt.Printf("Review %d: %s - %s\n", i, band, title)
  })
}

func main() {
  ExampleScrape()
}
```

## Related Projects

- [Goq][goq], an HTML deserialization and scraping library based on goquery and struct tags.
- [andybalholm/cascadia][cascadia], the CSS selector library used by goquery.
- [suntong/cascadia][cascadiacli], a command-line interface to the cascadia CSS selector library, useful to test selectors.
- [asciimoo/colly](https://github.com/asciimoo/colly), a lightning fast and elegant Scraping Framework
- [gnulnx/goperf](https://github.com/gnulnx/goperf), a website performance test tool that also fetches static assets.
- [MontFerret/ferret](https://github.com/MontFerret/ferret), declarative web scraping.

## Support

There are a number of ways you can support the project:

* Use it, star it, build something with it, spread the word!
  - If you do build something open-source or otherwise publicly-visible, let me know so I can add it to the [Related Projects](#related-projects) section!
* Raise issues to improve the project (note: doc typos and clarifications are issues too!)
  - Please search existing issues before opening a new one - it may have already been adressed.
* Pull requests: please discuss new code in an issue first, unless the fix is really trivial.
  - Make sure new code is tested.
  - Be mindful of existing code - PRs that break existing code have a high probability of being declined, unless it fixes a serious issue.

If you desperately want to send money my way, I have a BuyMeACoffee.com page:

<a href="https://www.buymeacoffee.com/mna" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;" ></a>

## License

The [BSD 3-Clause license][bsd], the same as the [Go language][golic]. Cascadia's license is [here][caslic].

[jquery]: http://jquery.com/
[go]: http://golang.org/
[cascadia]: https://github.com/andybalholm/cascadia
[cascadiacli]: https://github.com/suntong/cascadia
[bsd]: http://opensource.org/licenses/BSD-3-Clause
[golic]: http://golang.org/LICENSE
[caslic]: https://github.com/andybalholm/cascadia/blob/master/LICENSE
[doc]: http://godoc.org/github.com/PuerkitoBio/goquery
[index]: http://api.jquery.com/index/
[gonet]: https://github.com/golang/net/
[html]: http://godoc.org/golang.org/x/net/html
[wiki]: https://github.com/PuerkitoBio/goquery/wiki/Tips-and-tricks
[thatguystone]: https://github.com/thatguystone
[piotr]: https://github.com/piotrkowalczuk
[goq]: https://github.com/andrewstuart/goq
