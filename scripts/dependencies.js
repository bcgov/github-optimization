const { graphql } = require('@octokit/graphql');
require('dotenv').config();

const token = process.env.GITHUB_TOKEN;
const org = process.env.GITHUB_ORG;
const repo = process.env.GITHUB_REPOSITORY;

async function getRepositoryDependencyGraphManifests(org, repo) {
  const { repository } = await graphql(
    `
      {
        repository(name: "${repo}", owner: "${org}") {
          dependencyGraphManifests(first: 100){
            totalCount
            pageInfo {
              endCursor
              hasNextPage
            }
            nodes {
              id
              filename
              dependenciesCount
              dependencies (first: 100){
                totalCount 
              }
            }
          }
        }
      }
      `,
    {
      headers: {
        authorization: `token ${token}`,
        accept: 'application/vnd.github.hawkgirl-preview+json',
      },
    }
  );

  return repository;
}

async function main() {
  const repository = await getRepositoryDependencyGraphManifests(org, repo);
  console.log(repository);
}

main();
