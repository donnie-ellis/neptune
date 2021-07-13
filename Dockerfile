FROM golang:1.16.5

ENV GOOS=linux
ENV GOARCH=arm

RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build -o neptune .
CMD [ "/app/neptune" ]
