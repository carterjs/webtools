# `webtools`

`github.com/carterjs/webtools`

A collection of utilities for building web front ends with Go.

## `assets`

- Serve files with HTTP cache headers
- Minify assets with [tdewolff/minify](https://github.com/tdewolff/minify)
- Generate unique asset filenames to break caching on new releases

## `cache`

- Cache function executions
- Continue to serve stale results when functions fail

## `graphql`

- Make GraphQL requests 
- Parse responses into Go types using generics

## `templates`

- Recursively parse HTML template (`template/html`) files in a directory
- Minify templates with [tdewolff/minify](https://github.com/tdewolff/minify)
