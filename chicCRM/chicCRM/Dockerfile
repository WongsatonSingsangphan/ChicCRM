FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go get github.com/gin-gonic/gin/binding@v1.8.1
RUN go get go.mongodb.org/mongo-driver/x/mongo/driver@v1.11.1

RUN go build -o /chicCRM

EXPOSE 8888

CMD ["/chicCRM"]