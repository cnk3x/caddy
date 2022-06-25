#!/bin/bash

(curl -SsL https://caddyserver.com/api/packages | jq) >package.json
