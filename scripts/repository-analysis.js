const fs = require('fs');
const path = require('path');
const csv = require('fast-csv');
const _ = require('lodash');

const template = fs.readFileSync('repository-analysis-template.html');

const compiled = _.template(template);

const result = {
  primaryLanguages: {},
  allLanguages: {},
};

const schema = {
  repository: {},
  owner: {},
  name_with_owner: {},
  url: {},
  homepage_url: {
    type: String,
    count: `% of repositories has "Homepage Url"`,
  },
  description: {
    type: String,
    count: `% of repositories has "Description"`,
  },
  forks: {},
  is_fork: {
    type: Boolean,
    count: `% of forked repositories`,
  },
  is_mirror: {
    type: Boolean,
    count: `% of mirror repositories`,
  },
  is_private: {
    type: Boolean,
    count: `% of private repositories`,
  },
  is_archived: {
    type: Boolean,
    count: `% of archived repositories`,
  },
  is_template: {
    type: Boolean,
    count: `% of template repositories`,
  },
  stars: {},
  disk_usage: {},
  has_issues_enabled: {
    type: Boolean,
    count: `% of repositories with "Issues" enabled`,
  },
  has_projects_enabled: {
    type: Boolean,
    count: `% of repositories with "Projects" enabled`,
  },
  has_wiki_enabled: {
    type: Boolean,
    count: `% of repositories with "Wiki" enabled`,
  },
  merge_commit_allowed: {
    type: Boolean,
    count: `% of repositories that allows "Merge Commit"`,
  },
  rebase_merge_allowed: {
    type: Boolean,
    count: `% of repositories that allows "Rebase Merge"`,
  },
  squash_merge_allowed: {
    type: Boolean,
    count: `% of repositories that allows "Squash Merge"`,
  },
  created_at: {},
  updated_at: {},
  pushed_at: {},
  open_add_project_lifecycle_badge_issues: {
    type: Number,
    count: `% of repositories that have repoMountie issue "Lifecycle Badge"`,
  },
  open_lets_use_common_phrasing_issues: {
    type: Number,
    count: `% of repositories that have repoMountie issue "Common Phrasing"`,
  },
  open_add_missing_topics_issues: {
    type: Number,
    count: `% of repositories that have repoMountie issue "Missing Topics"`,
  },
  open_repomountie_issues: {
    type: Number,
    count: `% of repositories that have repoMountie issues`,
  },
  ministry_codes_count: {
    type: Number,
    count: `% of repositories that have Ministry codes`,
  },
  license_file_exists: {
    type: Boolean,
    count: `% of repositories that have "License"`,
  },
  readme_file_exists: {
    type: Boolean,
    count: `% of repositories that "README"`,
  },
  contributing_file_exists: {
    type: Boolean,
    count: `% of repositories that "Contributing"`,
  },
  code_of_conduct_file_exists: {
    type: Boolean,
    count: `% of repositories that "Code of Conduct"`,
  },
  changelog_file_exists: {
    type: Boolean,
    count: `% of repositories that "Change log"`,
  },
  security_file_exists: {
    type: Boolean,
    count: `% of repositories that "SECURITY"`,
  },
  support_file_exists: {
    type: Boolean,
    count: `% of repositories that "SUPPORT"`,
  },
  readme_references_license: {
    type: Boolean,
    count: `% of repositories that reference LICENSE in README`,
  },
  binaries_not_present: {
    type: Boolean,
    count: `% of repositories that contain Binaries`,
  },
  test_directory_exists: {
    type: Boolean,
    count: `% of repositories that contain Test Directory`,
  },
  integrates_with_ci: {
    type: Boolean,
    count: `% of repositories that integrates with CI`,
  },
  code_of_conduct_file_contains_email: {
    type: Boolean,
    count: `% of repositories that contain Email in Code of Conduct`,
  },
  source_license_headers_exist: {
    type: Boolean,
    count: `% of repositories that contain License sections in source codes`,
  },
  github_issue_template_exists: {
    type: Boolean,
    count: `% of repositories that contain Issue Templates`,
  },
  github_pull_request_template_exists: {
    type: Boolean,
    count: `% of repositories that contain Pull Request Templates`,
  },
  package_count: {
    type: Number,
    count: `% of repositories that contain Packages`,
  },
  project_count: {
    type: Number,
    count: `% of repositories that contain Projects`,
  },
  release_count: {
    type: Number,
    count: `% of repositories that contain Releases`,
  },
  submodule_count: {
    type: Number,
    count: `% of repositories that contain Sub modules`,
  },
  deploy_key_count: {
    type: Number,
    count: `% of repositories that contain Deploy Keys`,
  },
  topic_count: {
    type: Number,
    count: `% of repositories that contain Topics`,
  },
  license: {
    type: String,
    group: 'Licenses',
  },
  code_of_conduct: {
    type: String,
    group: 'Code Of Conduct',
  },
  days_open: {},
  issue_count: {},
  pr_count: {},
  commit_count: {},
  average_issue_count_per_day: {},
  average_pr_count_per_day: {},
  average_commit_count_per_day: {},
  default_branch_name: {
    type: String,
    group: 'Default Branches',
  },
  languages: {
    type: Array,
    custom: (data, row) => {
      const langs = data.split('_');
      const primaryLang = langs.length > 0 ? langs[0] : '';

      if (!result.primaryLanguages[primaryLang]) result.primaryLanguages[primaryLang] = 1;
      else result.primaryLanguages[primaryLang]++;

      langs.forEach((lang) => {
        if (!result.allLanguages[lang]) result.allLanguages[lang] = 1;
        else result.allLanguages[lang]++;
      });
    },
  },
  fork_pr_count: {
    type: Number,
    count: `% of repositories that contain Outside Contributions`,
  },
  review_count: {
    type: Number,
    count: `% of repositories that contain Pull Request Reviews`,
  },
};

async function main() {
  _.each(schema, (val, key) => {
    if (val.count) result[key] = 0;
    if (val.group) result[_.camelCase(`${key}_group`)] = {};
  });

  let totalCount = 0;

  fs.createReadStream(path.resolve(__dirname, '../notebook/dat/bcgov', 'master.csv'))
    .pipe(csv.parse({ headers: true }))
    // pipe the parsed input into a csv formatter
    .pipe(csv.format({ headers: true }))
    // Using the transform function from the formatting stream
    .transform((row, next) => {
      totalCount++;

      Object.keys(row).forEach((field) => {
        if (!schema[field]) return;

        if (schema[field].count) {
          switch (schema[field].type) {
            case String:
              if (row[field].length > 0) result[field]++;
              break;
            case Boolean:
              if (row[field] === 'True') result[field]++;
              break;
            case Number:
              if (Number(row[field]) > 0) result[field]++;
              break;
            default:
              break;
          }
        }

        if (schema[field].group) {
          const key = _.camelCase(`${field}_group`);

          if (!result[key][row[field]]) result[key][row[field]] = 1;
          else result[key][row[field]]++;
        }

        const customProcess = schema[field].custom || _.noop;
        customProcess(row[field], row);
      });

      next(null);
    })
    .on('finish', () => {
      const toPerc = (val) => ((val / totalCount) * 100).toFixed(2);
      const toArr = (obj) => {
        let j = 1;
        return _.reduce(
          obj,
          (ret, val, key) => {
            ret.push({ key: j, name: key, value: val, perc: toPerc(val) });
            j++;
            return ret;
          },
          []
        );
      };

      const counts = [];
      const groups = [];

      let i = 1;
      _.each(schema, (val, key) => {
        if (val.count) {
          counts.push({ key: i, name: val.count, value: result[key], perc: toPerc(result[key]) });
          i++;
        }

        if (val.group) {
          const name = _.camelCase(`${key}_group`);
          const group = toArr(result[name]);
          groups.push({ name: val.group, value: group });
        }
      });

      groups.unshift({ name: 'General Stastics', value: counts });
      groups.push({ name: 'Primary Languages', value: toArr(result.primaryLanguages) });
      groups.push({ name: 'All Languages', value: toArr(result.allLanguages) });

      result.totalCount = totalCount;
      result.groups = groups;

      fs.writeFile(path.resolve(__dirname, '../docs', 'index.html'), compiled(result), 'utf8', () =>
        process.exit()
      );
    });
}

main();
