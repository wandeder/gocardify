FROM golang:alpine
LABEL authors="Dmirtii Morozov"

WORKDIR /app

COPY ../../anki_bot /app
RUN go build -o main .
ENTRYPOINT ["./main"]