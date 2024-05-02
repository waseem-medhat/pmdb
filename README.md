# PMDb (work in progress)

![PMDb banner](https://i.imgur.com/DB0RFdF.png)

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
- [Templ](https://templ.guide/quick-start/installation)
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

- Compile the Tailwind CSS by running the `tw-build` NPM script. If you want to
run the Tailwind watcher, you'll have to run the `tw-watch` script in a
separate terminal.

```bash
npm run tw-build
```

- Run `air` (at the project root) which will compile and auto-reload the Go
server and Templ components.

```bash
air
```

## Tech Stack

### Technologies:

- **Go**: backend
- **Templ**: templating
- **TailwindCSS**: styling Go templates
- **HTMX**: UI interactivity
- **Turso**: database
- **Goose**: SQL migration tool
- **SQLC**: SQL-to-Go code generation tool
- **Air**: live reloading

### External resources:

- [**Icônes**](https://icones.js.org/): SVG icons
- [**TMDB**](https://developer.themoviedb.org/docs/getting-started): API for
movie information

### Architectural Design

The backend is built as a set of loosely coupled services as displayed in this chart.

![chart](https://i.imgur.com/CrMGM6H.png)

<!---
```
flowchart TD
    fe(["Frontend"]) --- movies
    fe --- home
    fe --- nowplaying

    subgraph Services
    nowplaying --- tmdbapi
    
    home --- auth
    home --- tmdbapi

    movies --- auth
    movies --- tmdbapi
    end
    
    auth --- db(["DB layer (SQLC-generated)"])
    tmdbapi --- api(["fa:fa-network-wired TMDB API"])
    home --- db
    movies --- db
    db --- database(["fa:fa-database Turso Database"])
```
-->

## Contributing

This project is still a work in progress and mainly for learning purposes, so
_for now_ I only accept issues for pointing out bugs or feature suggestions.
Thanks!
