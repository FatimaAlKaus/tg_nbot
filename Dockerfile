FROM golang:alpine as builder
ARG ssh_prv_key
RUN apk add git openssh
WORKDIR /go/src/app
COPY . .
# authorize SSH Host
RUN mkdir -p /root/.ssh && \
    chmod 0700 /root/.ssh && \
    ssh-keyscan github.com > /root/.ssh/known_hosts
# Add the keys and set permissions
RUN echo "$ssh_prv_key" > /root/.ssh/id_rsa && \
    chmod 600 /root/.ssh/id_rsa
# change urls
RUN git config --global url.ssh://git@github.com.insteadOf https://github.com
# download packages
RUN go mod tidy
#build
RUN go build -o /go/bin/app cmd/bot/main.go

#final stage
FROM alpine:latest
RUN mkdir bot && mkdir config
COPY --from=builder /go/bin/app /bot/
COPY --from=builder /go/src/app/config/dev.yml /config/dev.yml
ENTRYPOINT /bot/app
EXPOSE 80
