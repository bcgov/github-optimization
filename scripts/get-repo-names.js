const { graphql } = require("@octokit/graphql");
require("dotenv").config();
const fs = require("fs");
const token = process.env.GITHUB_TOKEN

async function getAllRepositoryNamesByOrg(org){
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
    repoNames = repoNames.concat(nodes.map(node => node.name))
  } while (hasNextPage)
  return repoNames;
}

async function main(){
  const stream = fs.createWriteStream("./data/repoNames.json", { flags: "a" });
  stream.write(
    "[ \n"
  );
  const repoNames = await getAllRepositoryNamesByOrg("bcgov")
  repoNames.forEach(repoName => {
    stream.write(`"${repoName}", \n`)
  })
  stream.write(
    "] \n"
  );
  stream.end();
}

main()
