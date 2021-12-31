# go-identity

<a href="https://github.com/greenpau/go-identity/actions/" target="_blank"><img src="https://github.com/greenpau/go-identity/workflows/build/badge.svg?branch=master"></a>
<a href="https://pkg.go.dev/github.com/greenpau/go-identity" target="_blank"><img src="https://img.shields.io/badge/godoc-reference-blue.svg"></a>

## Overview

This purpose is a library for the managing user identities for applications.
The core of the library is `User` and `Database` data structures.

The `User` data structure captures the dynamics of user identities in the
United States.

The `Database` data strcuture allows managing these identities. Currently,
the `Database` provides a way of managing local users for a web application.

The key concurrency features of the `Database` are:

* Only one Go routine is allowed making modifications to users at a time.
  During that time, the entire database locks.
* Keeps user identities in `Users` slice of the `Database` data
  structure. The elements of the slice are pointers of `User` data structure
  The slice only grows in size.
* Keeps references to user identities in a number hashes for faster lookup.
  The keys in the hashes are strings and the value is either a single
  pointer to `User` or a slice of pointers to `User` instances. If a reference
  keeps unique values, then it is a single pointer, e.g. username. Otherwise,
  e.g. in the case of being a part of a company, it is a slice.

The following keys are unique across the database:

* ID
* Username
* EmailAddress: a user can have multiple emails, but the emails must
  be unique across the database.
