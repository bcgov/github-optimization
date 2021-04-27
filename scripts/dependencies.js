const { graphql } = require("@octokit/graphql");
require("dotenv").config();
const fs = require("fs");

const token = process.env.GITHUB_TOKEN;

async function getAllRepositoryTopicsByOrg(org, repo){
  let hasNextPage;
  let after = "";
  let repoNames = [];

  // do {
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
          accept: 'application/vnd.github.hawkgirl-preview+json'
        },
      }
    );
    // const { endCursor } = pageInfo;
    console.log(JSON.stringify(repository, null, 2));
    // hasNextPage = pageInfo.hasNextPage;
    // if (hasNextPage) after = endCursor;
    // repoNames = repoNames.concat(nodes.map(node => {
    //   return {repoName: node.name, license: node.licenseInfo ? node.licenseInfo.name : "Null"}
    // }))
  // } while (false)
  return repoNames;
}

async function main(){
  // const stream = fs.createWriteStream("./data/licenses.csv", { flags: "a" });
  // stream.write(
  //   "Repository, license \n"
  // );
  const repoTopics = await getAllRepositoryTopicsByOrg("bcgov", "bc-sans")
  // repoTopics.forEach(repo => {
  //   const {repoName, license} = repo;
  //   stream.write(`${repoName}, ${license} \n`)
  // })
  // stream.end();
}

main()
