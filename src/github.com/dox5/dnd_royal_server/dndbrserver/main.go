package main

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "os"
  "strconv"

  "github.com/dox5/dnd_royal_server/api"
  "github.com/dox5/dnd_royal_server/model"
)

type FormValueGetter interface {
    FormValue(key string) string
}

type APIRequestHandler func(FormValueGetter) (interface{}, error)

type RegisterEndpoint func(method string,
                           uri string,
                           callback APIRequestHandler)

func handleApiRequest(handleRequest APIRequestHandler,
                      writer http.ResponseWriter,
                      request *http.Request,
                      method string) {

    if request.Method != method {
        msg := fmt.Sprintf("Incorrect method for endpoint. Expected %s but got %v",
                           method,
                           request.Method)
        http.Error(writer, msg, http.StatusMethodNotAllowed)
        return
    }

    response, err := handleRequest(request)

    if err != nil {
        msg := fmt.Sprintf("Request Handling failed: %s", err)
        http.Error(writer, msg, http.StatusInternalServerError)
        return
    }

    if response != nil {
        formattedResponse, err := json.Marshal(response)

        if err != nil {
            msg := fmt.Sprintf("Failed to format JSON response: %s", err)
            http.Error(writer, msg, http.StatusInternalServerError)
            return
        }

        writer.Header().Add("Content-Type", "application/json")
        writer.Write(formattedResponse)
    }
}

func mkHandlerFunc(method string, handler APIRequestHandler) http.HandlerFunc {
    return func(writer http.ResponseWriter, request *http.Request) {
        handleApiRequest(handler, writer, request, method)
    }
}

func roomIdentifierFromForm(parameters FormValueGetter) (model.Identifier, error) {
    roomIdString := parameters.FormValue("RoomId")

    if len(roomIdString) == 0 {
      return 0, fmt.Errorf("Must provide RoomId URL paramter")
    }

    roomId, err := strconv.ParseUint(roomIdString, 10, 64)

    if err != nil {
      return 0, fmt.Errorf("Room ID must be uint64: %s", err)
    }

    return model.Identifier(roomId), nil
}

func createUserAPIHandler(logger *log.Logger, rooms *RoomManager) http.Handler {
    apiMux := http.NewServeMux()

    endponts := []struct {method string;
                          url string;
                          callback APIRequestHandler} {

        {http.MethodPost,
         "/create",
         func(FormValueGetter) (interface{}, error) {
             room := rooms.Create()
             response := api.RoomCreateResponse{
                 RoomId: room.Id(),
                 GameMasterId: room.GameMaster().Id() }

             logger.Printf("Created room: %+v\n", response)

             return response, nil
        }},

        {http.MethodGet,
         "/map",
         func(parameters FormValueGetter) (interface{}, error) {
            _, err := roomIdentifierFromForm(parameters)

            if err != nil {
                return nil, err
            }

            return nil, fmt.Errorf("Not implemented")
        }}}

    for _, endpoint := range endponts {
        apiMux.HandleFunc(endpoint.url,
                          mkHandlerFunc(endpoint.method, endpoint.callback))
    }

    return apiMux
}

func main() {
    logger := log.New(os.Stdout, "", log.LUTC | log.Ldate | log.Ltime)
    logger.Println("~~ Starting DND Battle Royal Server ~~")
    mux := http.NewServeMux()

    rooms := NewRoomManager()
    fogController := NewFogEndpoint(logger, rooms)



    mux.Handle("/api/v1/room/",
               http.StripPrefix("/api/v1/room",
                                createUserAPIHandler(logger, rooms)))

    mux.Handle("/api/v1/fog/", http.StripPrefix("/api/v1/fog", fogController))

    mux.Handle("/api/v1/token/",
               http.StripPrefix("/api/v1/token",
                                MakePlayerTokenEndpoint(rooms, logger)))

    mux.Handle("/", http.FileServer(http.Dir("/web_static")))

    http.ListenAndServe(":8000", mux)
}
