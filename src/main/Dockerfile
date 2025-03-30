FROM golang:1.21.4

ENV CGO_ENABLED=0

RUN apt-get update && apt-get install -y netcat-traditional

RUN mkdir /code
WORKDIR /code

COPY runserver.sh /
RUN chmod +x /runserver.sh

COPY runfeeder.sh /
RUN chmod +x /runfeeder.sh

EXPOSE 80

COPY . /code/

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz -C /usr/local/bin

RUN go install github.com/spf13/cobra-cli@latest

RUN go build -o /app ./main.go

COPY entrypoint.sh /
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
