const chalk = require('chalk')

const { listEnvironments } = require('@keystone.sh/core/lib/commands/list')

const { CommandSignedIn } = require('../../lib/commands')

class ListCommand extends CommandSignedIn {
  async run() {
    await this.withUserSession(async userSession => {
      const project = await this.getProjectName()

      await listEnvironments(userSession, { project })
    })
  }
}

ListCommand.description = `Lists environments
`

ListCommand.examples = [chalk.blue('$ ks list environments')]

module.exports = ListCommand
