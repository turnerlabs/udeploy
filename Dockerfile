FROM golang:1.13

ARG version=0.0.0-unknown.0

COPY . /deploy/  

WORKDIR /deploy

# auto version dependencies to force browser cache refresh
RUN sed -i'.av' -E "s/\?auto-version=[0-9]+/\?auto-version=$RANDOM/g" vue/pages/**/*.html
RUN find vue/pages -name '*.av' -delete

RUN go build -ldflags '-X main.version='$version

CMD ./udeploy
