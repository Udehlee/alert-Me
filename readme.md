## alert-Me
Alert Me is a price drop alert system tracks submitted E-commerce products and notify users when there is a reduction in the price.

### How it works
A user submits a product URL (Jumia or konga), which is then published to a RabbitMQ queue. A consumer service called PriceCheck picks up the URL, scrapes the product's name and price, and stores the details in the database. Periodically, another service called SendForRecheck retrieves the saved product data, scrapes the current price again, and compares it with the database stored value. If a yhere is price drop, the system triggers a notification, which is published to the price_drop_alert queue.

### Technologies Used

- Go (Gin) 
- Postgres
- RabbitMQ
- Docker and docker-compose

### Setup and installation

- Clone the repository:

```sh 
git clone https://github.com/Udehlee/alert-Me.git
```
```sh
cd alert-Me
 ```
- Install dependencies 
```sh
go mod tidy
```

- Create .env file and fill it with your credentials as shown in the .env.example
- In the  .env.example
```sh
 NAME_SELECTOR
PRICE_SELECTOR
```
should either be jumia or Konga for now.

- Start the application with
 ```sh
 docker-compose up --build
```
The sever is listening on http://localhost:8000
the rabbitMQ will be listening http://localhost:15672, login as
 ```sh
name: guest
password: guest
```

 starting the application will apply  the following migration files in the internals/db/migrations folder and create:

- a products table that will hold all the details of the scraped product_url

### Api Endpoints


```sh
POST /submit
```