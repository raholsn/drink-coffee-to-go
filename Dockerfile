FROM golang

ARG app_env
ENV APP_ENV $app_env

COPY ./app /go/src/github.com/user/myProject/app
WORKDIR /go/src/github.com/user/myProject/app

RUN go get ./
RUN go build

CMD ["app"]
	
EXPOSE 8000