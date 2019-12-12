const {Command} = require('@oclif/command')
const {createUserSession} = require('../lib/blockstackLoader')
const {read} = require('../lib/cliStorage')
const chalk = require('chalk')

class WhoAmICommand extends Command {
  async run() {
    const session = await read({path: this.config.configDir + '/', filename: "session.json"})
    const userSession = createUserSession(session)
    if(userSession && userSession.isUserSignedIn()){
      const userData = userSession.loadUserData()
      this.log(`${chalk.bgHex('#f56565')(' ')} ${chalk.hex('#f56565').bold(` keystone. `)}`)
      this.log(`\n▻ You are connected under ${chalk.bold(userData.username)}`)
      this.log(`▻ You can logout with: ${chalk.yellow(`$ ks logout`)}`)

    }
    
  }
}

WhoAmICommand.description = `Shows the blockstack id of the currently logged in user
`

WhoAmICommand.examples = [
  chalk.yellow('$ ks whoami')
]

module.exports = WhoAmICommand

