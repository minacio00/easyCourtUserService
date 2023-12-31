FROM golang:alpine

WORKDIR /app
COPY . . 
RUN go mod download
COPY . .
RUN go build -o main .
EXPOSE 8081
ENV ADDR=8081
CMD [ "./main" ]