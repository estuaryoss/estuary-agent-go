FROM meterup/ubuntu-golang:1.15

ENV APP_DIR /app
ENV PORT 8080
# ENV HTTP_AUTH_TOKEN None
# ENV COMMAND_TIMEOUT 1200

RUN mkdir $APP_DIR
WORKDIR $APP_DIR

COPY estuary-agent-go $APP_DIR
ADD https://estuary-agent-go.s3.eu-central-1.amazonaws.com/4.1.0/runcmd-linux $APP_DIR/runcmd
ADD https $APP_DIR/https

RUN chmod +x estuary-agent-go
RUN chmod +x runcmd

CMD ["./estuary-agent-go"]
