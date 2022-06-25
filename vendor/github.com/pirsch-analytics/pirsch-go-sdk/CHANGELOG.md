# Changelog

## 1.7.1

* added single access token that don't require to query an access token using oAuth
* updated dependencies

## 1.7.0

* added more options to hit and event options

## 1.6.1

* fixed conversion goal return type

## 1.6.0

* added hourly visitor statistics, listing events, os and browser versions
* added filter options
* updated return types with new and modified fields
* improved error messages

## 1.5.1

* added endpoint for total visitor statistics

## 1.5.0

* added endpoint to extend sessions
* added entry page statistics
* added exit page statistics
* added number of sessions to referrer statistics
* added city statistics
* added entry page, exit page, city, and referrer name to filter

## 1.4.0

* added method to send events
* added reading event statistics
* fixed filter parameters to read statistics

## 1.3.0

* added `source` and `utm_source` to referrers
* added methods to read statistics
* updated dependencies

## 1.2.0

* added screen width and height to `Hit`
* improved refresh mechanism with wait time and ignoring obsolete requests

## 1.1.3

* fixed refreshing token on first request

## 1.1.2

* fixed refreshing token more often than needed

## 1.1.1

* added missing DNT (do not track) header

## 1.1.0

* removed deprecated package io/ioutil, the minimum Go version is now 1.16

## 1.0.0

Initial release.

## 0.4

* fixed 502 error and refreshing token
* added logger to `ClientConfig`

## 0.3

* hack to get around 502 responses for now

## 0.2

* fixed reading referrer from request

## 0.1

This is the first release for the Pirsch beta. The first version only includes sending hits to Pirsch. We will keep adding functionality.
