FROM golang:latest
WORKDIR /app
COPY go.mod go.sum  ./
RUN go mod tidy
COPY . .
COPY gcp_creds.json ./gcp_creds.json
ENV GOOGLE_APPLICATION_CREDENTIALS=/app/gcp_creds.json
RUN go build -o main ./predict_micro.go
CMD [ "/app/main" ]