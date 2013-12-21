#!/bin/bash

set -e

PORT=3001
FEEDS="http://dave.cheney.net/category/golang/feed/atom http://blog.gopheracademy.com/feed.atom http://feeds.feedburner.com/zen20 http://blog.campoy.cat/feeds/posts/default http://blog.golang.org/feed.atom http://www.goinggo.net/feeds/posts/default http://dominik.honnef.co/atom.xml"

go build github.com/davecheney/planetgolang
env PORT=$PORT ./planetgolang -template=$PWD/templates -static=$PWD/static $FEEDS
