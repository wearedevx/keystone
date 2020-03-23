const { flags } = require('@oclif/command')
// const chalk = require('chalk')
const chalk = require('chalk')

const {
  listAllMembers,
  listEnvMembers,
} = require('@keystone.sh/core/lib/commands/list')

const { CommandSignedIn } = require('../../lib/commands')

class ListCommand extends CommandSignedIn {
  async run() {
    await this.withUserSession(async userSession => {
      const env = await this.getProjectEnv()
      const project = await this.getProjectName()

      if (flags.project) {
        await listAllMembers(userSession, { project })
      } else {
        await listEnvMembers(userSession, {
          project,
          env,
          isProjectMembers: false,
        })
      }
    })
  }
}

ListCommand.description = `list members from current environment or project`

ListCommand.examples = [chalk.blue('$ ks member list ')]

ListCommand.flags = {
  ...CommandSignedIn.flags,
}

module.exports = ListCommand
