# horinezumi
Random wiki framework


# Running

## Environment variables
You need a .env file with 2 variables:
    ADDR: The web address to listen and serve
    PAGEDATAURL: The secret url to your database.

## Dependencies
Make sure to get the dependencies first:
* ![Chi - The HTTP library](github.com/go-chi/chi/v5)
* ![Pgx - The Postgres database driver](github.com/jackc/pgx/v5)
* ![Godotenv - Used to load the .env file](github.com/joho/godotenv)
* ![Go-JWT - The JWT library](github.com/golang-jwt/jwt/v5)


To run with make:
`make run`

To build the project
`make` or `make build`

I debug run with VScode