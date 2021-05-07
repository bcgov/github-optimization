# Github Optimization

[![Lifecycle:Experimental](https://img.shields.io/badge/Lifecycle-Experimental-339999)](Redirect-URL)

The goal of this GitHub Optimization Project is to provide tangible progress and recommendations for BC Government that will nurture and grow a healthy, compliant, and vibrant open-source development community. ​


## Features

This repository includes the following artifacts developed to understand the repositories in `bcgov` and `BCDevOps` GitHub organizations:

- `grafana-github-plugin`: a forked repository of the official grafana [`github-datasource`](https://github.com/grafana/github-datasource)
- `scripts-go`: a set of golang scripts to collect repository data from GitHub API
- `scripts`: a set of javascript scripts to collect repository data from GitHub API
- `notebook`: [`Jupyter`](https://jupyter.org/) notebook to merge repository data files into the master csv file

## Usage

- `grafana-github-plugin`: For information on setting up the grafana dashboard, please refer to its [README](/grafana-github-plugin/README.md).
- `scripts-go`: To run the go scripts, there are two variables expected, `GITHUB_TOKEN` and `GITHUB_ORG`. The token requires a personal access token
from your github account, and the github organization can be set to an existing organization, e.g. bcgov. To set the tokens, run:

  - `export GITHUB_TOKEN=<your PAT>`
  - `export GITHUB_ORG=<your organization>`

- `scripts`: before running the scripts, follow the .env.example file to add a personal access token.
- `notebook`: run `jupyter notebook` from the notebook directory to launch the notebooks.

## Requirements

This project was setup using the tool [asdf](https://asdf-vm.com/#/), and required tool version can be setup using it if you have asdf installed.
If not, the required tools and versions are located in `.tool-versions` and can be installed.

## Project Status

This project is in the experimental stage. 

## Goals/Roadmap

Please see our public project [backlog](https://github.com/bcgov/github-optimization/projects/2) for more information about the roadmap and goals.

## Getting Help or Reporting an Issue

If you have difficulties or problems, please open an issue and we will be happy to get back to you.

## How to Contribute

If you would like to contribute to the guide, please see our [CONTRIBUTING](CONTRIBUTING.md) guideleines.
"Please note that this project is released with a [Contributor Code of Conduct](CODE_OF_CONDUCT.md). By participating in this project you agree to abide by its terms."
