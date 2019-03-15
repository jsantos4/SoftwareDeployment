FROM golang:latest

ENV APP /home/jsantos4/go/src/SoftwareDeployment
WORKDIR /home/jsantos4/go/src/SoftwareDeployment

ADD . $APP

ENV LOGGLY_TOKEN fd30ee90-c128-4f11-b5af-2b00fc3a2afc

RUN cd ${APP} && go get -v github.com/jamespearly/loggly && go get -v github.com/robfig/cron
RUN go build main.go

CMD ./main

