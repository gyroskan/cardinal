# Cardinal

The api used to interact with cardinal database used for a discord bot.

You can find an exemple of the current api at [cardinal.gyroskan.com](https://cardinal.gyroskan.com/swagger/index.html)

## Usage

### Build from source

```sh
go mod download
go build -o cardinal
./cardinal
```

Do not forget to provide a .env file containing:

- SECRET
- DB_USER
- DB_HOST
- DB_PWD
- DB_NAME

### Docker-compose

```yaml
version: "3.7"
services:
  cardinal_api:
    image: gyroskan/cardinal_api:latest
    container_name: cardinal_api
    restart: unless-stopped
    environment:
      - SECRET=test-secret_used#for!token/sign
      - DB_HOST=db:3306
      - DB_USER=root
      - DB_PWD=root
      - DB_NAME=cardinal
    depends_on:
      - "db"
    ports:
      - "5005:5005"
    networks:
      - cardinal
      - db
  db:
    container_name: cardinal_api_db
    image: mariadb
    restart: unless-stopped
    environment:
      MARIADB_ROOT_PASSWORD: root
      MARIADB_DATABASE: cardinal
    volumes:
      - ./data:/var/lib/mysql
    networks:
      - db

networks:
  db:
  cardinal:
    external: true
```

-------

### Migrations

WIP

-------

The api is now available on the port 5005.  
You can see the swagger documentation
at <localhost:5005/swagger/index.html>  
and the api is available at <localhost:5005/api/v1>
