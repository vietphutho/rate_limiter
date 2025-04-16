#  API Rate Limiter (Golang)

Dự án này triển khai một hệ thống API Rate Limiting bằng Golang, hỗ trợ nhiều thuật toán kiểm soát tốc độ truy cập như:

 **Fixed Window**

 **Sliding Window Log**

 **Token Bucket**

 **Redis-backed Token Bucket** (cho distributed system)

 **Redis Pub/Sub** đồng bộ token giữa các instance

##  Mục Tiêu Dự Án
Thiết kế một middleware RateLimiter giúp:

- Giới hạn số request từ mỗi user trong khoảng thời gian nhất định
- Hỗ trợ nhiều chiến lược giới hạn tốc độ
- Dễ dàng mở rộng cho môi trường distributed (dùng Redis)
- Có thể dùng như middleware trong các HTTP API

##  Cấu Trúc Thư Mục
```plaintext
.
├── limiter/
│   ├── fixed_window.go
│   ├── sliding_window.go
│   ├── token_bucket.go
│   ├── redis_token_bucket.go
│   ├── redis_pubsub.go
│   ├── middleware.go
│   └── limiter_test.go
├── main.go
├── go.mod
└── README.md



