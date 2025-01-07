module github.com/justyntemme/gobolt

go 1.23.4

replace github.com/justyntemme/gobolt/server => /Users/justyntemme/Documents/code/gobolt/server/

replace github.com/justyntemme/gobolt/dom => /Users/justyntemme/Documents/code/gobolt/dom

require (
	github.com/gomarkdown/markdown v0.0.0-20241205020045-f7e15b2f3e62
	github.com/sirupsen/logrus v1.9.3
)

require golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
