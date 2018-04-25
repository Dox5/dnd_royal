package api

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strconv"

    "github.com/dox5/dnd_royal_server/model"
)

func RoomIdFromRequest(request *http.Request) (model.Identifier, error) {
    roomIdString := request.FormValue("RoomId")

    if len(roomIdString) == 0 {
      return 0, fmt.Errorf("Missing RoomId parameter")
    }

    roomId, err := strconv.ParseUint(roomIdString, 10, 64)

    if err != nil {
      return 0, fmt.Errorf("RoomId must be uint64: %s", err)
    }

    return model.Identifier(roomId), nil
}

func FormatResponse(writer http.ResponseWriter,
                    response interface{},
                    err error) {
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

func ParseJsonRequest(request *http.Request, v interface{}) error {
    buf := make([]byte, request.ContentLength, request.ContentLength)
    n, err := io.ReadFull(request.Body, buf)

    if int64(n) != request.ContentLength {
        return fmt.Errorf("Expcted to read %d bytes but actually read %d bytes",
                          request.ContentLength,
                          n)
    }

    err = json.Unmarshal(buf, v)
    return err
}
