const analysis = require('./repository-analysis');

analysis({ orgName: 'bcgov', source: 'bcgov/master.csv', outputFilename: 'index.html' });
