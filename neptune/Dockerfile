FROM golang:1.16.5

RUN mkdir /app
COPY ./neptune /app
WORKDIR /app
RUN go build -o trident .
CMD [ "/app/trident" ]
