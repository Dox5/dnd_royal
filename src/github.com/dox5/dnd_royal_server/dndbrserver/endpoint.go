package main

import (
    "fmt"
    "net/http"

    "github.com/dox5/dnd_royal_server/api"
)

type methodPath struct {
    method string
    path string
}

type EndpointHandler func(*http.Request) (interface{}, error)

type Endpoint struct {
    handlers map[methodPath]EndpointHandler
}

func NewEndpoint() *Endpoint {
    return &Endpoint{handlers: make(map[methodPath]EndpointHandler)}
}


// path is like /getMyThing
func (e *Endpoint) Register(path string,
                            method string,
                            handler EndpointHandler) {
    key := methodPath{method: method,
                      path: path}
    e.handlers[key] = handler
}

// Note: by-value receiver!
func (e Endpoint) ServeHTTP(writer http.ResponseWriter,
                            request *http.Request) {
    path := request.URL.Path
    var response interface{}
    var err error

    key := methodPath{method: request.Method,
                      path: path}
    handler, foundIt := e.handlers[key]

    if !foundIt {
        msg := fmt.Sprintf("No handler found for %+v", key)
        http.Error(writer, msg, http.StatusInternalServerError)
        return
    }

    response, err = handler(request)

    api.FormatResponse(writer, response, err)
}
