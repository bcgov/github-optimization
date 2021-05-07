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

const RANGE_COUNT = 7;

const schema = {
  repository: {},
  owner: {},
  name_with_owner: {},
  url: {},
  homepage_url: {
    type: String,
    count: `% of repositories that contain "Homepage Url"`,
  },
  description: {
    type: String,
    count: `% of repositories that contain "Description"`,
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
    count: `% of repositories that allow "Merge Commit"`,
  },
  rebase_merge_allowed: {
    type: Boolean,
    count: `% of repositories that allow "Rebase Merge"`,
  },
  squash_merge_allowed: {
    type: Boolean,
    count: `% of repositories that allow "Squash Merge"`,
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
    count: `% of repositories that contain "Ministry codes"`,
  },
  license_file_exists: {
    type: Boolean,
    count: `% of repositories that contain "License"`,
  },
  readme_file_exists: {
    type: Boolean,
    count: `% of repositories that contain "README"`,
  },
  contributing_file_exists: {
    type: Boolean,
    count: `% of repositories that contain "Contributing"`,
  },
  code_of_conduct_file_exists: {
    type: Boolean,
    count: `% of repositories that contain "Code of Conduct"`,
  },
  changelog_file_exists: {
    type: Boolean,
    count: `% of repositories that contain "Change log"`,
  },
  security_file_exists: {
    type: Boolean,
    count: `% of repositories that contain "SECURITY"`,
  },
  support_file_exists: {
    type: Boolean,
    count: `% of repositories that contain "SUPPORT"`,
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
    count: `% of repositories that contain "Issue Templates"`,
  },
  github_pull_request_template_exists: {
    type: Boolean,
    count: `% of repositories that contain "Pull Request Templates"`,
  },
  package_count: {
    type: Number,
    count: `% of repositories that contain "Packages"`,
  },
  project_count: {
    type: Number,
    count: `% of repositories that contain "Projects"`,
  },
  release_count: {
    type: Number,
    count: `% of repositories that contain "Releases"`,
  },
  submodule_count: {
    type: Number,
    count: `% of repositories that contain "Sub-modules"`,
  },
  deploy_key_count: {
    type: Number,
    count: `% of repositories that contain "Deploy Keys"`,
  },
  topic_count: {
    type: Number,
    count: `% of repositories that contain "Topics"`,
  },
  license: {
    type: String,
    group: 'Licenses',
  },
  code_of_conduct: {
    type: String,
    group: 'Code Of Conduct',
  },
  days_open: {
    type: Number,
    range: true,
  },
  issue_count: {
    type: Number,
    range: true,
  },
  pr_count: {
    type: Number,
    range: true,
  },
  commit_count: {
    type: Number,
    range: true,
  },
  branch_protection_rule_count: {
    type: Number,
    range: true,
  },
  avg_issue_count_per_day: {
    type: Number,
    range: true,
  },
  avg_pr_count_per_day: {
    type: Number,
    range: true,
  },
  avg_commit_count_per_day: {
    type: Number,
    range: true,
  },
  default_branch_name: {
    type: String,
    group: 'Default Branches',
  },
  languages: {
    type: Array,
    custom: (data = '', row) => {
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
    range: true,
  },
  review_count: {
    type: Number,
    count: `% of repositories that contain Pull Request Reviews`,
    range: true,
  },
};

async function main({ orgName, source, outputFilename }) {
  _.each(schema, (val, key) => {
    if (val.count) result[key] = 0;
    if (val.group) result[_.camelCase(`${key}_group`)] = {};
  });

  let totalCount = 0;
  const ranges = {};

  fs.createReadStream(path.resolve(__dirname, '../notebook/dat', source))
    .pipe(csv.parse({ headers: true }))
    // pipe the parsed input into a csv formatter
    .pipe(csv.format({ headers: true }))
    // Using the transform function from the formatting stream
    .transform((row, next) => {
      totalCount++;

      Object.keys(row).forEach((key) => {
        const field = _.snakeCase(key);

        if (!schema[field]) {
          if (totalCount === 1) console.log(`schema definition for ${field} is missing`);
          return;
        }

        const value = row[key];

        if (schema[field].count) {
          switch (schema[field].type) {
            case String:
              if (value.length > 0) result[field]++;
              break;
            case Boolean:
              if (value === 'True') result[field]++;
              break;
            case Number:
              if (Number(value) > 0) result[field]++;
              break;
            default:
              break;
          }
        }

        if (schema[field].group) {
          const key = _.camelCase(`${field}_group`);

          if (!result[key][value]) result[key][value] = 1;
          else result[key][value]++;
        }

        if (schema[field].range && schema[field].type === Number) {
          const key = _.camelCase(`${field}_range`);
``
          if (!ranges[key]) ranges[key] = [Number(value)];
          else ranges[key].push(Number(value));
        }

        const customProcess = schema[field].custom || _.noop;
        customProcess(value, row);
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
          groups.push({ name: val.group, value: group, chartType: 'horizontalBar' });
        }

        if (val.range) {
          const name = _.camelCase(`${key}_range`);
          const array = ranges[name];

          if (!array) return true;

          const min = _.min(array);
          const max = _.max(array);
          const range = (max - min) / RANGE_COUNT;

          if (min === 0 && max === 0) return true;
          if (range < 1) return true;

          const group = [];

          for (let x = 0; x < RANGE_COUNT; x++) {
            group.push({
              key: x,
              name: `${Math.floor(range * x)} - ${Math.floor(range * (x + 1) - 1)}`,
              value: 0,
              perc: 0,
            });
          }

          array.forEach((val) => {
            group[parseInt((val - min - 1) / range)].value++;
          });

          for (let x = 0; x < RANGE_COUNT; x++) {
            group[x].perc = toPerc(group[x].value);
          }

          groups.push({
            name: val.range === true ? _.startCase(key) : val.name,
            value: group,
            chartType: 'bar',
            sortBy: val.sortBy || 'key',
            sortOrder: val.sortOrder || 'asc',
          });
        }
      });

      groups.unshift({ name: 'General Stastics', value: counts, chartType: 'horizontalBar' });
      groups.push({ name: 'Primary Languages', value: toArr(result.primaryLanguages), chartType: 'horizontalBar' });
      groups.push({ name: 'All Languages', value: toArr(result.allLanguages), chartType: 'horizontalBar' });

      result.totalCount = totalCount;
      result.groups = groups;

      result.orgName = orgName;
      result.sourceUrl = `https://github.com/bcgov/github-optimization/blob/main/notebook/dat/${source}`;

      fs.writeFile(path.resolve(__dirname, '../docs', outputFilename), compiled(result), 'utf8', () => process.exit());
    });
}

module.exports = main;
