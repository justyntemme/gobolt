module github.com/justyntemme/gobolt

go 1.23.4

replace github.com/justyntemme/gobolt/server => /Users/justyntemme/Documents/code/gobolt/server/

replace github.com/justyntemme/gobolt/dom => /Users/justyntemme/Documents/code/gobolt/dom

replace github.com/justyntemme/gobolt/template => /Users/justyntemme/Documents/code/gobolt/template

require (
	github.com/gomarkdown/markdown v0.0.0-20241205020045-f7e15b2f3e62
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/text v0.21.0
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
)

require golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
