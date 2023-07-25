# <p align="center">Check-In 👉🔴</p>

![GitHub](https://img.shields.io/github/license/XDoubleU/check-in)
[![Build & Lint & Test](https://github.com/XDoubleU/check-in/actions/workflows/main.yml/badge.svg)](https://github.com/XDoubleU/check-in/actions/workflows/main.yml)
[![codecov](https://codecov.io/gh/XDoubleU/check-in/branch/main/graph/badge.svg?token=8IY0BGQ5RW)](https://codecov.io/gh/XDoubleU/check-in)
![GitHub last commit](https://img.shields.io/github/last-commit/XDoubleU/check-in)

This web application allows you to manage check-ins at multiple locations. Users can anonymously check-in by simply pressing a button. The app utilizes a websocket connection to display live location capacity on your website, allowing you to provide up-to-date information to your visitors.

Originally developed for Brugge Studentenstad, a non-profit organization that arranges student activities and manages study locations in Bruges, this app has been used to provide live capacity updates for each location on their website. Additionally, students are asked to select their school, providing Brugge Studentenstad with valuable data insights.

<p align="center">
   <img src="https://user-images.githubusercontent.com/54279069/232328182-92de6ebb-ce44-44c4-9796-6e6ef62fb7c6.jpg" style="height: 20em" />
   <br/>
   <em>Rubens styled painting of students, with their backpacks and books, lining up to press a red round button on a big screen when arriving at a library. - Generated using Bing Image Creator.</em>
</p>

## How to run locally?

1. Clone the repo
2. Start the web-client, API and database using `docker-compose up --build`
3. Apply migrations to database (in `api` dir) using `make db/migrations/up`
4. Create admin user (in `api` dir) using `make run/cli/createadmin u=admin p=admin`
5. Go to `http://localhost:3000` for the web-client and `http://localhost:8000` for the API

## Websocket integration

An example on how to use the websocket integration on your own website can be found [here](./integration/script.js).

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

For more information please read the [CONTRIBUTING](./CONTRIBUTING.md) document.

## License

[GPLv3](./LICENSE)
