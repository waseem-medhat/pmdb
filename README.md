# PMDb (work in progress)

![PMDb banner](https://i.imgur.com/TLD2Sty.png)

**Personal Movie Database**

PMDb is your personal space for movies! Here you can rate and review movies you
watched and make watchlists for movies you want to watch.

## Getting Started

### Installation

To get a local environment for this project, you need to have the following
installed:

- [Go](https://go.dev/dl/)
- Node.js (for TailwindCSS) either [directly](https://nodejs.org/en) or with
[NVM](https://github.com/nvm-sh/nvm)
- [TailwindCSS](https://tailwindcss.com/docs/installation)
- [Goose](https://github.com/pressly/goose?tab=readme-ov-file#install)
- [SQLC](https://docs.sqlc.dev/en/latest/overview/install.html)
- [Air](https://github.com/cosmtrek/air) (optional but highly recommended)

You also need to set up a [Turso database](https://docs.turso.tech/quickstart).

### Usage

After cloning the repo, setting up a database, and installing all dependencies:

- Copy the contents of `.env.example` into `.env` then supply `DBURL` and
`DBTOKEN` according to your database instance

```bash
cp ~/.env.example .env
```

- Execute the database migrations via Goose _in the `./sql/schema` directory_

```bash
goose turso "libsql://[DATABASE].turso.io?authToken=[TOKEN]" up
```

- Run `air` (at the project root) which will compile and auto-reload both the
Go server and TailwindCSS

```bash
air
```

## Tech Stack

- **Go**: backend + templating
- **TailwindCSS**: styling Go templates
- **HTMX**: UI interactivity
- **Turso**: database
- **Goose**: SQL migration tool
- **SQLC**: SQL-to-Go code generation tool
- **Air**: live reloading

## Contributing

This project is still a work in progress and mainly for learning purposes, so
_for now_ I only accept issues for pointing out bugs or feature suggestions.
Thanks!
