FROM golang:1.24-alpine AS build_stage
WORKDIR /consumer_app
COPY . .
RUN go mod tidy
RUN go build -o binary_app consumer.go  

FROM alpine AS run_stage
WORKDIR /app_binary
COPY --from=build_stage /consumer_app /app_binary/
# RUN chmod +x ./my_super_app

CMD [ "./binary_app" ]