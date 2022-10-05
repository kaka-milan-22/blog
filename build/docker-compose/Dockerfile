FROM golang:1.19.2-buster

WORKDIR /usr/local/app
COPY . .

RUN   apt update && apt install netcat -y && \
      go build -o ./bin/ 

EXPOSE 8888


CMD ["./build/docker-compose/wait-for","172.17.0.1:3306","--","./bin/blog"]