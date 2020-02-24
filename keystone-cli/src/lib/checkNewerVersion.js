const chalk = require('chalk')

const pkg = require('../../package')

const checkForUpdate = require('update-check')

module.exports = () =>
  checkForUpdate(pkg, {
    interval: 3600000, // For how long to cache latest version (default: 1 day)
  })
    .then(result => {
      if (result) {
        console.log(
          `${chalk.bgRed('UPDATE AVAILABLE')} version ${chalk.yellow(
            result.latest
          )} is out. Run ${chalk.blue(
            `npm i -g '@keystone.sh/cli@latest'`
          )} to install it.`
        )
      }
    })
    .catch(err => console.error(`Failed to check for updates: ${err}`))
