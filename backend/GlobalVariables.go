package main

import "encoding/json"

var null = []byte ("null")
var emptyArray = json.RawMessage("[]")

type VideoJSON struct {
	UserVideos json.RawMessage
	AllVideos json.RawMessage
}

type AudioJson struct {
	UserMusic json.RawMessage
	AllMusic json.RawMessage
}