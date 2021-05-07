const axios = require("axios");
const cheerio = require("cheerio");

const getHTML = (url) =>
  axios
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

const parsePage = async (html) => {
  const $ = cheerio.load(html);
  const boxes = $(".org-row .org");
  let repoCount = 0;

  for (let x = 0; x < boxes.length; x++) {
    const $$ = cheerio.load(cheerio.html(boxes[x]));
    const linkUrl = $$("a").attr("href");

    const prefix = "https://github.com/";

    const orgMetaUrl = `https://api.github.com/orgs/${linkUrl.substr(
      prefix.length
    )}`;

    const orgMeta = await getHTML(orgMetaUrl);

    repoCount = (repoCount + orgMeta && orgMeta.public_repos) || 0;
  }

  return { orgCount: boxes.length, repoCount };
};

const getGovernmentCount = async () => {
  const url = `https://government.github.com/community/`;

  const content = await getHTML(url);

  return parsePage(content);
};

async function main() {
  const governmentCount = await getGovernmentCount();
  console.log(governmentCount);
}

main();
