FROM node:20 as build-frontend
WORKDIR /build

COPY ./frontend/package.json ./frontend/yarn.lock .
RUN yarn install --frozen-lockfile

COPY ./frontend .
RUN yarn build

FROM golang:1.25 as build-server

WORKDIR /build

COPY go.mod go.sum .
RUN go mod download

COPY . .
COPY --from=build-frontend /build/dist ./frontend/dist

RUN CGO_ENABLED=0 go build -buildvcs=false -o ./bin/server ./cmd/server

FROM golang:1.25 as build-worker

WORKDIR /build

COPY go.mod go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -buildvcs=false -o ./bin/worker ./cmd/worker

FROM alpine:3.18

RUN apk add --no-cache ca-certificates

COPY --from=build-server /build/bin/server /usr/bin/server
COPY --from=build-worker /build/bin/worker /usr/bin/worker

EXPOSE 8080

ENV SERVER_PORT=8080

CMD ["/usr/bin/server"]
