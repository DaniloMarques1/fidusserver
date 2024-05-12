FROM golang:1.22.2

WORKDIR /fidus

COPY . ./
RUN go mod download
RUN go mod tidy

EXPOSE 8080

RUN go build

CMD ["./fidusserver"]
