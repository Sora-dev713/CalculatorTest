FROM golang:1.24-alpine AS build
WORKDIR /src/backend
COPY backend/go.mod ./
COPY backend ./
ARG SERVICE
RUN test -n "$SERVICE" && CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/service ./cmd/${SERVICE}

FROM alpine:3.22
RUN addgroup -S calculator && adduser -S calculator -G calculator
COPY --from=build /out/service /usr/local/bin/service
USER calculator
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/service"]
