FROM golang:1.23-alpine3.20 

RUN apk add --no-cache unzip curl gcc musl-dev libaio

RUN mkdir -p /opt/oracle/

RUN curl -o ic.zip https://download.oracle.com/otn_software/linux/instantclient/2326100/instantclient-basic-linux.x64-23.26.1.0.0.zip \
    && unzip ic.zip -d /opt/oracle/instantclient_23_26 \
    && rm ic.zip

ENV LD_LIBRARY_PATH="/opt/oracle/instantclient_23_26"
ENV CGO_ENABLED=1 

WORKDIR /usr/src/app

RUN go mod init dbmt
RUN go get github.com/godror/godror

COPY . .

RUN go build -o dbmt .

EXPOSE 5461

CMD ["./dbmt"]
