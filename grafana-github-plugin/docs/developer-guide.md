# Developer Guide

This is a very basic guide on how to set up your local environment, make the desired changes and see the result with a fresh Grafana Installation.

## Getting Started

Clone this repository into your local environment. The frontend code lives in the `src` folder, alongside the [plugin.json file](https://grafana.com/docs/grafana/latest/developers/plugins/metadata/). See [this grafana tutorial](https://grafana.com/docs/grafana/latest/developers/plugins/) to understand better how a plugin is structured and installed.

Backend code, written in Go, is located in the `pkg` folder.

## Requirements

For this standard execution, you will need the following tools:

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Yarn](https://classic.yarnpkg.com/en/docs/install)

## Running the development version

### Compiling the Backend

If you have made any changes to any `go` files, you can use [mage](https://github.com/magefile/mage) to recompile the plugin.

```sh
mage build:linux && mage reloadPlugin
```

## Workaround

These instructions are failing, but should work according to [grafana](https://github.com/grafana/grafana-plugin-sdk-go/tree/master/build), need to figure out why. 
In the meantime, make sure to get the following:
``` go
go get github.com/grafana/grafana-plugin-sdk-go/build
go get github.com/grafana/grafana-plugin-sdk-go/backend@v0.90.0
go get github.com/hashicorp/go-plugin@v1.2.2
go get github.com/grafana/grafana-plugin-sdk-go/data@v0.90.0
```

and run `mage buildAll` to compile the backend.

### Compiling the Frontend

After you made the desired changes, you can build and test the new version of the plugin using `yarn`:

```sh
yarn test # run all test cases
yarn dev # builds and puts the output at ./dist
```

Alternatively, you can have yarn watch for changes and automatically recompile them.

```sh
yarn watch
```

Now that you have a `./dist` folder, you are ready to run a fresh Grafana instance and put the new version of the datasource into [Grafana plugin folder](https://grafana.com/docs/grafana/latest/plugins/installation/).

### Docker Compose

We provide a [Docker Compose file](/docker-compose.yml) to help you to get started. When you call up `docker-compose up` inside the project folder, it will:

1. Run a new instance of Grafana from the master branch and map it into port `3090`.
1. Configure the instance to allow an unsigned version of `github-datasource` to be installed.
1. Map the current folder contents into `/var/lib/grafana/plugins`.

This is enough for you to see the Github Datasource in the datasource list at `http://localhost:3090/datasources/new`.

![Local Github Stats installation](./screenshots/local-plugin-install.png)

If you make further changes into the code, be sure to run `yarn dev` again and restart the Grafana instance.

## Create a pull request

After you are good to go, it is time to create a pull request to share your work with the community. Please read more about that [here](https://github.com/grafana/grafana/blob/master/contribute/create-pull-request.md).

## ASDF install

```sh
cat .tool-versions | cut -f 1 -d ' ' | xargs -n 1 asdf plugin-add
asdf plugin add mage https://github.com/ggilmore/asdf-mage.git
asdf plugin add oc https://github.com/sqtran/asdf-oc.git
asdf plugin-add docker-compose https://github.com/virtualstaticvoid/asdf-docker-compose.git
```

## Docker commands

```sh
yarn build
docker build -t grafana-ext .
docker run -d -p 3000:3000  grafana-ext
```
