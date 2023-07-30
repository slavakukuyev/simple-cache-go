﻿# simple-cache-go

This repository contains a simple implementation of an in-memory cache in Go (Golang). The cache supports setting key-value pairs with a time-to-live (TTL) expiration for entries. It uses a concurrent-safe design with a read-write mutex and provides a method to stop the cache from expiring entries before their TTL elapses.
