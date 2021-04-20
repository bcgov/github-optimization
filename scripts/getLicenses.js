const { graphql } = require("@octokit/graphql");
require("dotenv").config();
const fs = require("fs");

const token = process.env.GITHUB_TOKEN;

async function getAllRepositoryTopicsByOrg(org){
  let hasNextPage;
  let after = "";
  let repoNames = [];

  do {
    const { organization: {repositories: {pageInfo, nodes}} } = await graphql(
      `
      {
        organization(login: "${org}") {
          repositories(first: 100 ${after && `, after: "${after}"`}) {
            pageInfo {
              endCursor
              hasNextPage
            }
            nodes {
              name
              licenseInfo {
                name
              }
            }
          }
        }
      }
      `,
      {
        headers: {
          authorization: `token ${token}`,
        },
      }
    );
    const { endCursor } = pageInfo;
    hasNextPage = pageInfo.hasNextPage;
    if (hasNextPage) after = endCursor;
    repoNames = repoNames.concat(nodes.map(node => {
      return {repoName: node.name, license: node.licenseInfo ? node.licenseInfo.name : "Null"}
    }))
  } while (hasNextPage)
  return repoNames;
}

async function main(){
  const stream = fs.createWriteStream("./data/licenses.csv", { flags: "a" });
  stream.write(
    "Repository, license \n"
  );
  const repoTopics = await getAllRepositoryTopicsByOrg("bcgov")
  repoTopics.forEach(repo => {
    const {repoName, license} = repo;
    stream.write(`${repoName}, ${license} \n`)
  })
  stream.end();
}

main()
