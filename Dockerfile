FROM golang:1.15-alpine3.12

ENV APP_DIR /app
ENV PORT 8080
#ENV HTTP_AUTH_TOKEN None
#ENV COMMAND_TIMEOUT 1200

RUN mkdir $APP_DIR
WORKDIR $APP_DIR

COPY ./ $APP_DIR
ADD https://estuary-agent-go.s3.eu-central-1.amazonaws.com/4.1.0/runcmd-alpine $APP_DIR/runcmd

RUN chmod +x runcmd

CMD ["go", "run", "Main.go"]
