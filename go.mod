module github.com/ebookstore

go 1.16

require (
	github.com/aws/aws-sdk-go v1.48.13
	github.com/aws/aws-sdk-go-v2 v1.24.1
	github.com/aws/aws-sdk-go-v2/config v1.25.10
	github.com/aws/aws-sdk-go-v2/credentials v1.16.9
	github.com/aws/aws-sdk-go-v2/service/s3 v1.47.1
	github.com/aws/aws-sdk-go-v2/service/ses v1.19.6
	github.com/brianvoe/gofakeit/v6 v6.23.0
	github.com/bxcodec/faker/v3 v3.6.0
	github.com/gin-contrib/cors v1.4.0
	github.com/gin-contrib/size v0.0.0-20230212012657-e14a14094dc4
	github.com/gin-gonic/gin v1.8.1
	github.com/go-playground/validator/v10 v10.11.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/google/uuid v1.4.0
	github.com/jackc/pgx/v4 v4.13.0 // indirect
	github.com/lib/pq v1.10.2
	github.com/magiconair/properties v1.8.7
	github.com/redis/go-redis/v9 v9.4.0
	github.com/spf13/viper v1.8.1
	github.com/steinfletcher/apitest v1.5.15
	github.com/steinfletcher/apitest-jsonpath v1.7.2
	github.com/stretchr/testify v1.8.4
	github.com/stripe/stripe-go/v72 v72.67.0
	github.com/swaggo/files v0.0.0-20210815190702-a29dd2bc99b2
	github.com/swaggo/gin-swagger v1.3.1
	github.com/swaggo/swag v1.7.1
	github.com/testcontainers/testcontainers-go v0.27.0
	github.com/testcontainers/testcontainers-go/modules/localstack v0.27.0
	github.com/testcontainers/testcontainers-go/modules/redis v0.27.0
	github.com/ulule/limiter/v3 v3.10.0
	go.uber.org/zap v1.19.0
	golang.org/x/crypto v0.14.0
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.14
)
