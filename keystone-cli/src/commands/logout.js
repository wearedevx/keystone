const chalk = require('chalk')
const { SESSION_FILENAME } = require('@keystone.sh/core/lib/constants')
const { del } = require('../lib/cliStorage')
const { CommandSignedIn } = require('../lib/commands')

class LogoutCommand extends CommandSignedIn {
  async run() {
    await this.withUserSession(async userSession => {
      const userData = userSession.loadUserData()
      try {
        await del({
          path: `${this.config.configDir}/`,
          filename: SESSION_FILENAME,
        })
        this.log(
          `▻ Sign out from ${chalk.yellow(
            userData.username
          )} ${chalk.green.bold('✓')}`
        )
      } catch (error) {
        this.log(`▻ Can't log out: ${error.message}`)
      }
    })
  }
}

LogoutCommand.description = `Logs you out of your account and erase your session from this computer.
`

LogoutCommand.examples = [chalk.yellow('$ ks logout')]

module.exports = LogoutCommand
