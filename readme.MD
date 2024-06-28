# Synapsis

This is the backend for an example MVP online store application built with Go

## Features

- User registration and login
- View product list by category name
- Add products to the shopping cart
- View shopping cart items
- Delete items from the shopping cart
- Checkout and get invoice

## Tech Stack

- Go
- PostgreSQL
- Docker

## Getting Started

### Prerequisites

- Go 1.18+
- Docker

### Running the Application

1. Clone the repository:
    ```sh
    git clone https://github.com/armiariyan/synapsis.git
    cd synapsis
    ```

2. Build and run the Docker containers:
    ```sh
    docker-compose up --build
    ```

3. The application will be available at `http://localhost:3000`.

### Database Guide
1. Currently for example data in database, i manually use queries insert after the database and table has been created
2. The (optional) insert queries that i used:
    ```
    insert into product_categories (uuid, name)
values
    ('c54d278f-deec-4285-bb3a-c5cd95437099', 'electronics'),
    ('d7a271ee-bcbb-4d6d-a6fe-c4d0672304ef', 'books');
insert into products (uuid, name, category_id, price)
values
    ('401cfe92-eb8e-4326-a946-7707246061df', 'Iphone 15 Pro Max', 'c54d278f-deec-4285-bb3a-c5cd95437099', 10000000),
    ('9997e5dc-fcbe-49f4-89f2-5784ee2e8b71', 'Cosmos Wadesdes Fan', 'c54d278f-deec-4285-bb3a-c5cd95437099', 100000),
        ('5db2b2e0-dcc0-40e3-8720-315d06bdd464', 'Macbook M1 Pro', 'c54d278f-deec-4285-bb3a-c5cd95437099', 11000000),
        ('77c4514c-55f3-4209-8caa-3bc6dee0cabf', 'Book of Kedebook', 'd7a271ee-bcbb-4d6d-a6fe-c4d0672304ef', 150000),
    ('59e039c8-590a-411c-b362-377ccdf12be4', 'Bumi Manusia', 'd7a271ee-bcbb-4d6d-a6fe-c4d0672304ef', 100000);
    ```


### API Endpoints

- **Users**
  - `POST /v1/users/register`: Register a new user.
  - `POST /v1/users/login`: Login a user.
  - `POST /v1/users/checkout`: Checkout *all items* in carts.

- **Products**
  - `GET /v1/products`: List products, filter by category name.
  - `POST /v1/products/add-to-cart`: Add a product to user cart


- **Shopping Cart**
  - `GET /v1/carts`: View cart items.
  - `DELETE /carts/:itemUUID`: Delete item from cart.






