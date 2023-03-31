# go-oauth2-server

## Quick overview:
You can start by using the [Postman Collection](https://www.postman.com/qrestoque/workspace/pedro-tashima-s-public-workspace/collection/13233153-c52c7618-7e33-48ab-b855-f0b54e27134e?action=share&creator=13233153)    
> The requests need to be made in order, there are automatic tests that will set the collection variables for you.

* This Implementation is meant to be used by Public clients // TODO: add link to RFC 6749 2.1.  Client Types

## Disclaimers: 
* **NEVER** use this in production, this implementation **INSECURE**, but it's a good start point to understand how OAuth 2.0 works

## Table of Contents

1. [How to understand this repository](#how-to-understand-this-repository)
1. [Documentation](#documentation)
1. [Installation and Setup](#installation-and-setup)
1. [File Structure](#file-structure)

## How to understand this repository
1. First, you should have a basic understanding of how Oauth2 works.
1. Look at the [Postman Collection](https://www.postman.com/qrestoque/workspace/pedro-tashima-s-public-workspace/collection/13233153-c52c7618-7e33-48ab-b855-f0b54e27134e?action=share&creator=13233153), do some requests and try to understand how they compare to the sequence diagram below.
1. Go to the handlers folder, and see the handlers, try to undersand the rules.

## Documentation

This OAuth2.0 implementation is meant to be an example for Toolbox CloudPass. Follow the official documentation for more detailed instructions.    

[Toolbox Documentation](https://toolboxdigital.atlassian.net/wiki/spaces/DDP/pages/72293671/CloudPass+Integration+Guide+method+OAuth+2.0+Protocol)    
[Oauth Website](http://oauth.com)

> This repository implements `Note 2`, `Note 4` and `Note 6`    

![flow](https://toolboxdigital.atlassian.net/wiki/download/thumbnails/72293671/Flujo%20de%20autenticaci%C3%B3n%20Oauth2.jpg?version=1&modificationDate=1569931404030&cacheVersion=1&api=v2&width=1108&height=1921)

## Installation and Setup

1. Clone the repo: `git clone https://gitub.com/tashima42/go-oauth2-server.git`
1. Go to the directory: `cd go-oauth2-server`
1. Start the database: `docker-compose up -d`
1. Run the tests: `APP_DB_USERNAME=postgres APP_DB_PASSWORD=password APP_DB_NAME=postgres APP_PORT=8010 go test -v`
1. Run the app: `APP_DB_USERNAME=postgres APP_DB_PASSWORD=password APP_DB_NAME=postgres APP_PORT=8010 go run .`
