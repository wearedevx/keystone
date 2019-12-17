const checkout = require('./checkout')
const configure = require('./configure')

module.exports = {
  ...checkout,
  ...configure,
}
