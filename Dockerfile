FROM golang:1.17-alpine AS build
WORKDIR /src

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
# Generate doc with swaggo
RUN go get -u github.com/swaggo/swag/cmd/swag
RUN swag i -pd --generatedTime --overridesFile .swaggo -g api/api.go

RUN CGO_ENABLED=0 go build -o cardinal

FROM scratch
COPY --from=build /src/cardinal /cardinal
ENTRYPOINT ["/cardinal"]
