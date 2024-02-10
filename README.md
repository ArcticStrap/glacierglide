# horinezumi
A currently WIP backend wiki framework written in Go.

## Running

### Environment variables
You need a .env file with 3 variables:
    ADDR: The web address to listen and serve
    PAGEDATAURL: The secret url to your database.
    JWTCODE: The JWT secret.

Config variables:
    DEV: "1"=development mode ""=production mode

### Dependencies
* ![Chi](https://github.com/go-chi/chi) - The HTTP library
* ![Godotenv](https://github.com/joho/godotenv) - Used to load the .env file
* ![Go-JWT ](https://github.com/golang-jwt/jwt) - The JWT library

#### Database drivers
* ![Pgx](https://github.com/jackc/pgx/) - Postgres

I might plan to add MariaDB soon.

To get the Dependencies:
`make getdeps`

Generate the SSL certificates if running on HTTPS:
`make gencert`

To run with make:
`make run`

To build the project:
`make`

To clean the project:
`make clean`

To run unit tests:
`make test`

## Contributing

Open to contributions. Feel free to open up an issue or pull request.

## License

This project is licensed under the [MIT License](LICENSE)
