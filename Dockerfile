FROM golang:alpine AS build
RUN apk add --update --no-cache git
ADD . /src
RUN cd /src && CGO_ENABLED=0 go build -o sony-bravia-webui-remote

FROM scratch
COPY --from=build /src/sony-bravia-webui-remote /
COPY --from=build /src/ui /ui
EXPOSE 8080
ENV SONY_BRAVIA_IP "192.0.2.10"
ENV SONY_BRAVIA_PSK "0000"
ENTRYPOINT ["/sony-bravia-webui-remote"]
