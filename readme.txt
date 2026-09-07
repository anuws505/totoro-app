- init
go mod init totoro-app

- remove go.sum
rm go.sum

- clear cache ของโมดูล
go clean -cache -modcache

- จัดระเบียบ Module ใหม่ และโหลด Library ที่ขาด
go mod tidy

- อัพเดท library
go get -u ./...

- ติดตั้ง GORM (Core Package)
go get gorm.io/gorm

- ติดตั้ง Driver Postgres สำหรับ GORM
go get gorm.io/driver/postgres

- options
go get github.com/shopspring/decimal
go get github.com/google/uuid



- on local environment
docker compose up -d
docker compose down -v

DB_HOST=localhost \
DB_USER=user \
DB_PASS=password \
DB_NAME=totoro_db \
DB_PORT=5432 \
REDIS_ADDR=localhost:6379 \
AUTH_SECRET_KEY=a-string-secret-at-least-256-bits-long \
go run cmd/api/main.go

- docker build command
docker build -t totoro-app:v1 .

- docker option
docker image prune -f

- docker run command
docker run --network totoro-app_default \
  -p 8080:8080 \
  -e DB_HOST=db \
  -e DB_USER=user \
  -e DB_PASS=password \
  -e DB_NAME=totoro_db \
  -e DB_PORT=5432 \
  -e REDIS_ADDR=cache:6379 \
  -e AUTH_SECRET_KEY=a-string-secret-at-least-256-bits-long \
  totoro-app:v1

- or command with ".env" files
docker run --network totoro-app_default \
  -p 8080:8080 \
  -e DB_HOST=db \
  -e REDIS_ADDR=cache:6379 \
  --env-file .env \
  totoro-app:v1


curl -X POST http://localhost:8080/api/v1/auth/create \
-H "Content-Type: application/json" \
-d '{
  "mobile": "0812225656",
  "password": "1234"
}'

curl -X POST http://localhost:8080/api/v1/auth/login \
-H "Content-Type: application/json" \
-d '{
  "mobile": "0812225656",
  "password": "12345"
}'

curl -X POST http://localhost:8080/api/v1/auth/otp/request \
-d '{
  "mobile":"0812225656"
}'

curl -X POST http://localhost:8080/api/v1/auth/otp/verify \
-d '{
  "mobile":"0812225656",
  "otp":"410736"
}'

curl -X POST http://localhost:8080/api/v1/auth/logout \
-H "Authorization: Bearer vvvv"


curl -X PATCH http://localhost:8080/api/v1/user/profile \
-H "Authorization: Bearer vvvv" \
-H "Content-Type: application/json" \
-d '{
  "new_user_name": "NewName-5656"
}'


curl -X POST http://localhost:8080/api/v1/user/change-password \
-H "Authorization: Bearer vvvv" \
-H "Content-Type: application/json" \
-d '{
  "old_password": "1234",
  "new_password": "12345"
}'


curl -X PATCH http://localhost:8080/api/v1/admin/update-role \
-H "Authorization: Bearer vvvv" \
-H "Content-Type: application/json" \
-d '{
  "target_user_id": "z6XT2D_hzbkl",
  "new_role": "ADMIN"
}'



curl -X POST http://localhost:8080/api/v1/admin/result-master \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiM1FkZUdMQnA0V0d1IiwiZXhwIjoxNzc4OTA0MDE4LCJpYXQiOjE3Nzg4MTc2MTh9.EUPltHTRFCllSb1LwE2YVODbtAVoBNBfz92SV3Hmq44" \
-H "Content-Type: application/json" \
-d '{
  "date_reward": "2026-05-16 15:00:00",
  "results": {"DD": "10", "TTT": "210"}
}'


curl -X PATCH "http://localhost:8080/api/v1/admin/result-master/5" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiM1FkZUdMQnA0V0d1IiwiZXhwIjoxNzc4OTA0MDE4LCJpYXQiOjE3Nzg4MTc2MTh9.EUPltHTRFCllSb1LwE2YVODbtAVoBNBfz92SV3Hmq44" \
-H "Content-Type: application/json" \
-d '{
  "date_reward": "2026-05-15 19:00:00",
  "is_active": false,
  "results": { "DD": "10", "TTT": "160" }
}'



curl -X GET "http://localhost:8080/api/v1/items"
curl -X GET "http://localhost:8080/api/v1/items?with_deleted=true"


curl -X POST http://localhost:8080/api/v1/admin/item-master \
-H "Authorization: Bearer vvvv" \
-H "Content-Type: application/json" \
-d '{
  "item_sku": "ITEM-1331",
  "item_name": "Special 1331",
  "type": {
    "DIRECT": 25,
    "MIXED": 12
  },
  "code": "TTT",
  "code_value": "1331",
  "is_featured": false
}'


curl -X PATCH http://localhost:8080/api/v1/admin/item-master \
-H "Authorization: Bearer vvvv" \
-H "Content-Type: application/json" \
-d '{
  "item_sku": "ITEM-200",
  "item_name": "Special 200",
  "type": {
    "DIRECT": 50,
    "MIXED": 12
  },
  "code": "TTT",
  "code_value": "1331",
  "is_featured": true
}'


curl -X DELETE "http://localhost:8080/api/v1/admin/item-master?item_sku=ITEM-1331" \
-H "Authorization: Bearer vvvv"


curl -X POST "http://localhost:8080/api/v1/admin/item-master/restore?item_sku=ITEM-1331" \
-H "Authorization: Bearer vvvv"



curl -X GET "http://localhost:8080/api/v1/admin/promotions?with_deleted=true" \
-H "Authorization: Bearer vvvv"


curl -X POST http://localhost:8080/api/v1/admin/promotions \
-H "Authorization: Bearer vvvv" \
-H "Content-Type: application/json" \
-d '{
  "code": "SAVE250",
  "value": 250,
  "max_usage": 310,
  "expiry_date": "2026-12-31T23:59:59Z",
  "is_featured": false
}'


curl -X PATCH http://localhost:8080/api/v1/admin/promotions \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiM1FkZUdMQnA0V0d1IiwiZXhwIjoxNzc4ODI2NzIzLCJpYXQiOjE3Nzg3NDAzMjN9.Pm-xXKv-kqRJ3LyP1ZlH9xxqgk2F4_gQ7-4wuT6bVFQ" \
-H "Content-Type: application/json" \
-d '{
  "code": "SAVE30",
  "max_usage": 10,
  "is_featured": true
}'


curl -X DELETE "http://localhost:8080/api/v1/admin/promotions?code=SAVE25" \
-H "Authorization: Bearer vvvv"



curl -X POST http://localhost:8080/api/v1/sheets \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMzBpdWRTU1lJM3FjIiwiZXhwIjoxNzc4OTMwNDY1LCJpYXQiOjE3Nzg4NDQwNjV9.4eb-ud55Tfilz2NAq7Z_7PffH18QReI88B4W11RYHZs" \
-H "Content-Type: application/json" \
-d '{
  "items": [
    { "item_sku": "ITEM-10", "price": 10, "multiple": "DIRECT" },
    { "item_sku": "ITEM-200", "price": 40, "multiple": "MIXED" }
  ],
  "promotion_code": "SAVE20"
}'



{
  "Items": [
    {
      "item_sku": "ITEM-00",
      "item_name": "No.0",
      "price": 10,
      "type": {
        "DIRECT": 25,
        "MIXED": 4
      },
      "multiple": "DIRECT",
      "reward": 250,
      "prize": true,
      "code": "DD",
      "code_value": "00",
      "is_featured": false
    },
    {
      "item_sku": "ITEM-10",
      "item_name": "No.10",
      "price": 80,
      "type": {
        "DIRECT": 25,
        "MIXED": 4
      },
      "multiple": "DIRECT",
      "reward": 2000,
      "prize": false,
      "code": "DD",
      "code_value": "10",
      "is_featured": false
    },
    {
      "item_sku": "ITEM-200",
      "item_name": "No.200",
      "price": 40,
      "type": {
        "DIRECT": 50,
        "MIXED": 12
      },
      "multiple": "MIXED",
      "reward": 480,
      "prize": false,
      "code": "TTT",
      "code_value": "200",
      "is_featured": false
    }
  ],
  "Totals": {
    "Total": 130, // Items[0].Price + Items[1].Price + Items[2].Price
    "Discount": { "Code": "SAVE10", "Value": 10 },
    "GrandTotal": 120, // Totals.Total - Totals.Discount.Value
    "Reward": 250 // (Items[0].Prize(Items[0].Reward)) + (Items[1].Prize(Items[1].Reward)) + (Items[2].Prize(Items[2].Reward))
  },
  "ID": "ID"
  "SheetID": "SHEET-260511-xxx",
  "SheetCodeDD": "DD",
  "SheetCodeTTT": "TTT",
  "UserID": "xxxx",
  "UserName": "name-x",
  "CreatedAt": timestamp,
  "UpdatedAt": timestamp,
  "ExpiresAt": timestamp,
  "Status": "PAID",
  "PaidAt": timestamp,
  "RewardAt": "define_date_time",
  "RewardCompletedAt": timestamp,
  "PaymentCompletedAt": timestamp,
  "CompletedAt": timestamp,
  "AgentID": "yyyy",
  "AgentName": "name-y"
}

{
  "DateReward": "2026-05-16 15:00:00" // refer thai date or custom input
  "Type": {
    "DD": 10,
    "TTT": 160
  }
  "IsActive: true
}
