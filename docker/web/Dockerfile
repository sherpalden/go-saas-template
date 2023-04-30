FROM golang:1.20.0-alpine
RUN apk add inotify-tools
RUN go install github.com/rubenv/sql-migrate/...@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy 
COPY . ./

CMD sh /app/docker/web/run.sh