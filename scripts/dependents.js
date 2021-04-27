const fs = require("fs");
const { graphql } = require("@octokit/graphql");
const axios = require("axios");
const cheerio = require("cheerio");
const csv = require("fast-csv");

require("dotenv").config();

const token = process.env.GITHUB_TOKEN;
const org = process.env.GITHUB_ORG;

const getRepositoryNames = async (org) => {
  let hasNextPage;
  let after = "";
  let repoNames = [];

  do {
    const {
      organization: {
        repositories: { pageInfo, nodes },
      },
    } = await graphql(
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
    repoNames = repoNames.concat(nodes.map((node) => node.name));

    const { endCursor } = pageInfo;
    hasNextPage = pageInfo.hasNextPage;
    after = hasNextPage ? endCursor : null;
  } while (!!after);

  return repoNames;
};

const parseContents = (html) => {
  const $ = cheerio.load(html);
  const boxes = $("div.Box-row");

  const results = [];

  for (let x = 0; x < boxes.length; x++) {
    const $$ = cheerio.load(cheerio.html(boxes[x]));
    const orgrepo = cheerio.text($$('[data-repository-hovercards-enabled=""]'));
    const nospace = orgrepo.replace(/^\s+|\s+|\r+$/g, "");

    results.push(nospace);
  }

  const nextPage = $(".paginate-container a:contains('Next')").attr("href");

  return [results, nextPage];
};

const getDependents = async (org, repo) => {
  let url = `https://github.com/${org}/${repo}/network/dependents`;
  const dependents = [];

  console.log(`querying '${org}/${repo}'`);

  do {
    const content = await axios
      .get(url)
      .then((res) => res.data)
      .catch(function (error) {
        if (error.response) {
          console.log(error.response.data);
          console.log(error.response.status);
          console.log(error.response.headers);
        } else if (error.request) {
          console.log(error.request);
        } else {
          console.log("Error", error.message);
        }

        return null;
      });

    if (!content) break;

    const [results, nextPage] = parseContents(content);

    dependents.push(...results);

    console.log("nextPage", nextPage);
    url = nextPage;
  } while (!!url);

  return dependents;
};

async function main() {
  if (!org || !token) {
    console.error(`GITHUB_TOKEN and GITHUB_ORG must be defined.`);
  }
  const stream = fs.createWriteStream(`./data/dependents-${org}.csv`, {
    flags: "w",
  });

  const csvStream = csv.format({ headers: true });
  csvStream.pipe(stream).on("end", () => process.exit());

  const repoNames = await getRepositoryNames(org);
  console.log(`${repoNames.length} repositories found.`);

  for (let x = 0; x < repoNames.length; x++) {
    const dependents = await getDependents(org, repoNames[x]);

    csvStream.write({
      Repository: repoNames[x],
      Dependents: dependents.join(";"),
      "Dependent Count": dependents.length,
    });
  }

  csvStream.end();
}

main();
