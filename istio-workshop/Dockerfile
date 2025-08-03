FROM golang:1.19.6-buster AS build
WORKDIR /app
COPY . .
RUN go build -o app main.go

FROM gcr.io/distroless/base-debian11
WORKDIR /root/
COPY --from=build /app/app .
EXPOSE 8080
CMD ["./app"]