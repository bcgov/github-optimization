const util = require("util");
const exec = util.promisify(require("child_process").exec);
const fs = require('fs')

const org = "bcgov";

const defaults = [
  "license-file-exists",
  "readme-file-exists",
  "contributing-file-exists",
  "code-of-conduct-file-exists",
  "changelog-file-exists",
  "security-file-exists",
  "support-file-exists",
  "readme-references-license",
  "binaries-not-present",
  "test-directory-exists",
  "integrates-with-ci",
  "code-of-conduct-file-contains-email",
  "source-license-headers-exist",
  "github-issue-template-exists",
  "github-pull-request-template-exists",
  "javascript-package-metadata-exists",
  "ruby-package-metadata-exists",
  "java-package-metadata-exists",
  "python-package-metadata-exists",
  "objective-c-package-metadata-exists",
  "swift-package-metadata-exists",
  "erlang-package-metadata-exists",
  "elixir-package-metadata-exists",
  "license-detectable-by-licensee",
  "notice-file-exists"
]

const getRepoUrl = (org, repo) => {
  return `https://github.com/${org}/${repo}.git`;
};

function printToFile(repo, answer, stream) {
  return new Promise((resolve, reject) => {
    try {
      if (!answer || !answer.stdout) {
        throw new Error;
      }
      const results = JSON.parse(answer.stdout).results;
      stream.write(`${repo}`);
  
      results.forEach(result => {
        stream.write(`,${result.status}`)
      });
      stream.write(`\n`)
      resolve();
    }
    catch (e) {
      reject(e)
    }
  })
}


const lintRepos = async (org, repos) => {
  const stream = fs.createWriteStream("./dat/repolint-results.csv", { flags: "a" });
  stream.write('Repository')
  defaults.forEach(header => {
    stream.write(`, ${header}`);
  });
  stream.write(`\n`);
  for (let i = 0; i < repos.length; i++) {
    const repo = repos[i]
    console.log(`linting ${repo}, on ${i}th iteration`);
    try {
      const answer = await exec(`repolinter lint -g ${getRepoUrl(org, repo)} --format json || true`);
      await printToFile(repo, answer, stream);
    } catch (e) {
      console.log('huge err', e);
    }
  }
  stream.end();
};

let rawdata = fs.readFileSync('repoNames.json');
let parsedData = JSON.parse(rawdata);

lintRepos(org, parsedData);
