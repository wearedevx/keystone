const environments = require('./environments')
const members = require('./members')
const files = require('./files')
const projects = require('./projects')

module.exports = {
  ...environments,
  ...members,
  ...files,
  ...projects,
}
