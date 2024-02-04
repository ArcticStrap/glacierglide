# horinezumi
A currently WIP wiki software written in go.

## Running

### Environment variables
You need a .env file with 3 variables:
    ADDR: The web address to listen and serve
    PAGEDATAURL: The secret url to your database.
    JWTCODE: The JWT secret.

### Dependencies
Make sure to get the dependencies first:
* ![Chi](https://github.com/go-chi/chi) - The HTTP library
* ![Godotenv](https://github.com/joho/godotenv) - Used to load the .env file
* ![Go-JWT ](https://github.com/golang-jwt/jwt) - The JWT library

#### Database drivers
* ![Pgx](https://github.com/jackc/pgx/) - Postgres

I might plan to add MariaDB soon.

Generate the SSL certificates before running:
`make gencert`

To run with make:
`make run`

To build the project:
`make`

To clean the project:
`make clean`

To run unit tests:
`make test`

I debug run with VScode
