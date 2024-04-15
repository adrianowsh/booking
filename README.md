# Hotel reservation backend with Golang

## Project outline
- users   -> book room from an hotem
- admins  -> going to check reservation/bookings
- Auth    -> JWT tokens
- Hotels  -> CRUD API -> JSON
- Rooms   -> CRUD API -> JSON
- Scripts -> database -> seeding and migrations

## Resources
### MongoDb Driver
Documentation
```
https://mongodb.com/docs/drivers/go/current/quick-start

```

### gofiber
Documentation
```
https://gofiber.com
```

Installing gofiber
```
go get github.com/gofiber/fiber/v2
```

## Docker
### Instaling mongodb as a Docker container
```
docker run --name mondogb -d mongo:latest -p 27017:27017
```
