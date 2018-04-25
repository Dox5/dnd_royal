from golang:1-alpine

WORKDIR /go/
COPY src src/

RUN go get -d -v ./...
RUN go test ./...
RUN go install -v ./...

WORKDIR /
COPY web_static/ web_static/


CMD ["dndbrserver"]
