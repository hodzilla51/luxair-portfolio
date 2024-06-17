FROM golang:1.21.5

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /sample-backend-go ./cmd/sample-backend-go

EXPOSE 5000

ENV AWS_PROFILE=sample
ENV AWS_SDK_LOAD_CONFIG=1

CMD [ "/sample-backend-go" ]
