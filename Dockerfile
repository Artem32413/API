FROM golang:latest

WORKDIR /newappapi

COPY ./ ./
RUN go build -o main .
EXPOSE 8080
CMD [ "./main" ]