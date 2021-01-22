FROM docker.io/golang:1.15.6-alpine3.12 as build-stage
WORKDIR /app
COPY . .
RUN go build -o server ./cmd/main.go

FROm docker.io/alpine:3.12.3 as production-stage
COPY --from=build-stage /app/server /usr/local/bin/server
EXPOSE 3000
CMD ["server"]
