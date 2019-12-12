const chalk = require('chalk')
const { CommandSignedIn } = require('../lib/commands')

class WhoAmICommand extends CommandSignedIn {
  async run() {
    await this.withUserSession(async userSession => {
      const userData = userSession.loadUserData()
      this.log(`▻ You are connected under ${chalk.bold(userData.username)}`)
      this.log(`▻ You can logout with: ${chalk.yellow(`$ ks logout`)}`)
    })
  }
}

WhoAmICommand.description = `Shows the blockstack id of the currently logged in user
`

WhoAmICommand.examples = [chalk.yellow('$ ks whoami')]

module.exports = WhoAmICommand
