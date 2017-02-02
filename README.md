# Glide Pin

[![Build Status](https://travis-ci.org/multiplay/glide-pin.svg)](https://travis-ci.org/multiplay/glide-pin)
[![Go Report Card](https://goreportcard.com/badge/github.com/multiplay/glide-pin)](https://goreportcard.com/report/github.com/multiplay/glide-pin)

A plugin for glide that converts the lock file into the yaml file. Why might
you want this? Well, if you want a nice easy way to pin all the versions of all
the dependencies and sub-dependencies you have, this will take that info that
already exists in your glide.lock file and make it an explicit requirement in
your yaml file.

By default, the original yaml file has the dependencies in it replaced with the
values in the lock file, but the version strings kept from the original YAML.
You can pass the `--exact` flag to force the YAML file to have all version
strings replaced with the exact commit hash/version.

Loosely based on the code from [glide-cleanup][1] - thanks!

## Install

    go get github.com/multiplay/glide-pin

## Run

    glide pin

[1]:https://github.com/ngdinhtoan/glide-cleanup
