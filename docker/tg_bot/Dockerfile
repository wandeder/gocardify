FROM golang:latest as builder
LABEL authors="Dmirtii Morozov"
WORKDIR /app
COPY ./tg_bot /app
RUN go build -o main ./main/main.go
RUN chmod +x ./main/main

FROM golang:latest
WORKDIR /app
COPY --from=builder /app/main/main .
ENTRYPOINT ["./main"]