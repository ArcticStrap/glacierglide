# horinezumi
Random wiki framework


## Running

### Environment variables
You need a .env file with 2 variables:
    ADDR: The web address to listen and serve
    PAGEDATAURL: The secret url to your database.

### Dependencies
Make sure to get the dependencies first:
* ![Chi](github.com/go-chi/chi/v5) - The HTTP library
* ![Pgx](github.com/jackc/pgx/v5) - The Postgres database driver
* ![Godotenv](github.com/joho/godotenv) - Used to load the .env file
* ![Go-JWT ](github.com/golang-jwt/jwt/v5) - The JWT library


To run with make:
`make run`

To build the project
`make` or `make build`

I debug run with VScode