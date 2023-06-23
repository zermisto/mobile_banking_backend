FROM golang AS builder
WORKDIR /root/app
RUN apt update
RUN apt install bc

COPY go.mod go.mod
RUN go mod tidy
RUN go get github.com/prisma/prisma-client-go

COPY prisma prisma
RUN go run github.com/prisma/prisma-client-go generate
COPY . .
RUN go get -u

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./build/app
RUN echo "Size of binary = $(echo "$(stat -c%s '/root/app/build/app')/1024000" | bc -l) MB."

FROM frolvlad/alpine-glibc
RUN addgroup --system --gid 1001 nonroot
RUN adduser --system --uid 1001 nonroot

RUN apk add dumb-init
COPY --from=builder /root/app/build /
USER nonroot
ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/app"]