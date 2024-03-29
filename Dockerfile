# Build the Go API
FROM golang:latest AS builder
ADD . /app
WORKDIR /app/server
RUN go mod download
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /main .

# Build the React app
FROM node:alpine AS node_builder
COPY --from=builder /app/client ./
RUN npm install
RUN npm run build

# Final stage build
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Install Python and create a virtual environment
RUN apk add --no-cache python3 py3-pip
# Install netcat
RUN apk add --no-cache netcat-openbsd

# Setup venv
RUN python3 -m venv /venv
ENV PATH="/venv/bin:$PATH"
# Install Python Dependencies for ETL
COPY --from=builder /app/etl ./
RUN pip3 install --no-cache-dir -r requirements.txt

COPY --from=builder /main ./
COPY --from=node_builder /build ./web
RUN chmod +x ./main
EXPOSE 8080

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

CMD ["/entrypoint.sh"]