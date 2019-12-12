const add = require('./add')
const del = require('./delete')
const env = require('./env')
const init = require('./init')
const invite = require('./invite')
const list = require('./list')
const project = require('./project')
const pull = require('./pull')
const push = require('./push')

module.exports = {
  add,
  delete: del,
  env,
  init,
  invite,
  list,
  project,
  pull,
  push,
}
