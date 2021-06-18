FROM nunchistudio/blacksmith-enterprise:0.18.0-alpine

ADD ./ /fragment
WORKDIR /fragment

RUN rm -rf go.sum
RUN go mod tidy

EXPOSE 9090 9091
