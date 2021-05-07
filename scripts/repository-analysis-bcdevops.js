const analysis = require('./repository-analysis');

analysis({ orgName: 'bcdevops', source: 'bcdevops/master.csv', outputFilename: 'bcdevops.html' });
