const { graphql } = require("@octokit/graphql");
const { Octokit } = require("@octokit/rest");
require("dotenv").config();
const fs = require("fs");

let rawdata = fs.readFileSync('./data/repoNames.json');
let repoNames = JSON.parse(rawdata);

const token = process.env.GITHUB_TOKEN;
const octokit = new Octokit({
  auth: token,
  userAgent: 'myApp v1.2.3',
})

async function getAllRepositoryContributors(owner, repo){
  const contributors = await octokit.rest.repos.listContributors({
    owner,
    repo,
  });
  const contributorIds = contributors.data.map(contributor => contributor.node_id);
  return contributorIds
}

async function getAllRepositoryCollaboratorsByOrg(org, repo){
  let after = ""
  let hasNextPage = ""
  let collaboratorsWithDirectAccess = []

  do {
    const { organization: {repository: {collaborators}} } = await graphql(
      `
      {
        organization(login: "${org}") {
          repository(name: "${repo}") {
            collaborators (first: 100 ${after && ` after: "${after}"`}) {
              edges {
                permission
                node {
                  id
                }
              }
              pageInfo {
                hasNextPage
                endCursor
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
    const { pageInfo, edges } = collaborators;
    hasNextPage = pageInfo.hasNextPage;
    if (hasNextPage) after = pageInfo.endCursor;
    collaboratorsWithDirectAccess = collaboratorsWithDirectAccess.concat(edges.filter(edge => edge.permission !== "READ"))
  } while (hasNextPage)
  return collaboratorsWithDirectAccess.map(collaborator => collaborator.node.id);
}

async function main(){
  const failures = []
  const stream = fs.createWriteStream("./data/collaborators2.csv", { flags: "a" });
  stream.write(
    "Repository, Outside Contributor Count \n"
  );
  for (repoName of repoNames){
    // Running too many API requests concurrently triggers abuse mechanism, so grabbing one repo at a time
    try {
      console.log(`Grabbing data for repo ${repoName}`)
      const [collaboratorIds, contributorIds] = await Promise.all([getAllRepositoryCollaboratorsByOrg("bcgov", repoName), getAllRepositoryContributors("bcgov", repoName)]);
      const outsideContributorsCount = contributorIds.filter(id => !collaboratorIds.includes(id)).length;
      stream.write(`${repoName}, ${outsideContributorsCount} \n`)
    } catch (e) {
      console.log(`failed for repo ${repoName}`);
      failures.push(repoName);
    }
  }
  stream.end();
  console.log(failures)
}

main()
