# horinezumi
Random wiki framework


## Running

### Environment variables
You need a .env file with 2 variables:
    ADDR: The web address to listen and serve
    PAGEDATAURL: The secret url to your database.

### Dependencies
Make sure to get the dependencies first:
* ![Chi](https://github.com/go-chi/chi) - The HTTP library
* ![Pgx](https://github.com/jackc/pgx/) - The Postgres database driver
* ![Godotenv](https://github.com/joho/godotenv) - Used to load the .env file
* ![Go-JWT ](https://github.com/golang-jwt/jwt) - The JWT library


To run with make:
`make run`

To build the project
`make` or `make build`

I debug run with VScode