# 1 шаг
FROM golang:1.24-alpine AS build_stage
WORKDIR /my_super_app
COPY . .
RUN go mod tidy
RUN go build -o binary_app cmd/main.go  

# 2 шаг
FROM alpine AS run_stage
WORKDIR /app_binary
COPY --from=build_stage /my_super_app /app_binary/
# RUN chmod +x ./my_super_app
EXPOSE 8002/tcp
CMD [ "./binary_app"]