FROM golang:1.21.4

ENV CGO_ENABLED=0

RUN apt-get update && apt-get install -y netcat-traditional

RUN mkdir /code
WORKDIR /code

COPY runserver.sh /
RUN chmod +x /runserver.sh

EXPOSE 80

COPY . /code/

COPY entrypoint.sh /
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
