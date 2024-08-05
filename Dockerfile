FROM golang:1.22-alpine
WORKDIR /usr/src/app
COPY . .
RUN go build
EXPOSE 8090
CMD [ "./h1emu-charts"]
