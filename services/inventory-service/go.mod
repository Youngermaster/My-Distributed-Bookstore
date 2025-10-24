module github.com/youngermaster/distributed-bookstore/inventory-service

go 1.21

require (
	github.com/gofiber/fiber/v2 v2.51.0
	gorm.io/gorm v1.25.5
	gorm.io/driver/postgres v1.5.4
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/redis/go-redis/v9 v9.3.0
	google.golang.org/grpc v1.59.0
	google.golang.org/protobuf v1.31.0
	github.com/streadway/amqp v1.1.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	github.com/prometheus/client_golang v1.17.0
	github.com/google/uuid v1.4.0
)
