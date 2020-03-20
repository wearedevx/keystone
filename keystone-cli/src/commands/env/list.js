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

ListCommand.description = `List environments`

ListCommand.examples = [chalk.blue('$ ks env list')]

module.exports = ListCommand
