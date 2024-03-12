# GlacierGlide

![Penguin knowledge](./docs/static/banner.png)

GlacierGlide is an prototype-stage backend wiki framework written in Go. Basic features are still in development.

---

## Running

### Environment variables
You need a .env file with 3 variables:
    ADDR: The web address to listen and serve
    PAGEDATAURL: The secret url to your database.
    JWTCODE: The JWT secret.

Config variables:
```bash
DEV: Run in development or produciton mode: Set to 1 for the former.
WITHPOLARP: Use the PolarPages frontend. Set 1 to enable.
```

### Dependencies
* ![Chi](https://github.com/go-chi/chi) - The HTTP library
* ![Go-JWT ](https://github.com/golang-jwt/jwt) - The JWT library

#### Database drivers
* ![Pgx](https://github.com/jackc/pgx/) - PostgreSQL

To get the Dependencies:
`make getdeps`

Generate the SSL certificates if running on HTTPS:
`make gencert`

To generate the database schemas:
`make genschema`

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

This project is licensed under the [BSD 3 Clause License](LICENSE)
