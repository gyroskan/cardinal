FROM golang:1.16-alpine AS build
WORKDIR /src

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o cardinal

FROM scratch
COPY --from=build /src/cardinal /cardinal
ENTRYPOINT ["/cardinal"]
