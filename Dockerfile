FROM golang:1.20-alpine3.17
WORKDIR /apps
COPY go.mod ./
COPY go.sum ./
COPY ./main.go /apps/
COPY ./pkg/ /apps/pkg/
COPY ./version/ /apps/version/
COPY ./web/ /apps/web/


RUN go mod tidy

COPY *.go ./

RUN go build -o /nostr-client


CMD [ "/nostr-client" ]