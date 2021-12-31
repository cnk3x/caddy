# gopkg for Caddy v2

This is a port of the gopkg plugin to Caddy v2.

The plugin allows to easily setup vanity gopkg/go module names through a Caddy-powered website.

# Usage

`gopkg [path] [repo-uri]`

`gopkg [path] [vcs] [repo-uri]`

```
zikes.me {
  // default vcs is git
  gopkg /multistatus https://github.com/zikes/multistatus
  gopkg /chrisify https://github.com/zikes/chrisify

  // use mercurial
  gopkg /myrepo hg https://bitbucket.org/zikes/myrepo
}
```

The above would make the repos go get-able via `go get zikes.me/chrisify`,
`go get zikes.me/myrepo`, and `go get zikes.me/multistatus`.

If the urls are visited normally the browser will be redirected to the repo uri.

Once implemented, `go get` can enforce your import paths with
[import path checking](https://golang.org/cmd/go/#hdr-Import_path_checking).
