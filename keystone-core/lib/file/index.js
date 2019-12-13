const disk = require('./disk')
const gaia = require('./gaia')

module.exports = {
  ...disk,
  ...gaia,
}
