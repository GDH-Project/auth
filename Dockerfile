LABEL org.opencontainers.image.source = "https://github.com/GDH-Project/auth"

FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o main ./cmd

# ==================
# 메인 이미지
# ==================
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

ENV TZ=Asia/Seoul

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 50501

CMD ["./main"]