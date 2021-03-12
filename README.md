[![Coverage Status](https://coveralls.io/repos/github/xplorfin/docker-utils/badge.svg?branch=master)](https://coveralls.io/github/xplorfin/docker-utils?branch=master)
[![Renovate enabled](https://img.shields.io/badge/renovate-enabled-brightgreen.svg)](https://app.renovatebot.com/dashboard#github/xplorfin/docker-utils)
[![Build status](https://github.com/xplorfin/docker-utils/workflows/test/badge.svg)](https://github.com/xplorfin/docker-utils/actions?query=workflow%3Atest)
[![Build status](https://github.com/xplorfin/docker-utils/workflows/goreleaser/badge.svg)](https://github.com/xplorfin/docker-utils/actions?query=workflow%3Agoreleaser)
[![](https://godoc.org/github.com/xplorfin/docker-utils?status.svg)](https://godoc.org/github.com/xplorfin/docker-utils)
[![Go Report Card](https://goreportcard.com/badge/github.com/xplorfin/docker-utils)](https://goreportcard.com/report/github.com/xplorfin/docker-utils)

# What is this?

This is a helper library for interacting with docker by [entropy](https://entropy.rocks), particularly in continuous integration workflows. It provides a canonical way to create volumes, containers, images, etc and tear them down.
