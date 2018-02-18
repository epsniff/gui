# gui

## Status

A work in progress, when I get some freetime I'll chip a way at it and make incrmential improvments.  Currently all the work has been going into the infastruture.  Mostly I'm building out the work scheduler and communiication layers.  See: https://github.com/epsniff/gui/tree/master/src/server

## Goals

1. To create a pure Go alternative to ElasticSearch. 
2. To provider a simple query DSL using FilterQL https://github.com/araddon/qlbridge/blob/master/FilterQL.md.
3. Utlize Bleve's new Scorch index struture: https://github.com/blevesearch/bleve/blob/master/index/scorch/README.md.
4. Use a binary protocal like gRPC for node to node comunication, to remove the cost of using json/http.
5. Utlize a smart client simlar to how Google's Cloud Bigtable does, to remove the need for a cordinator node.
