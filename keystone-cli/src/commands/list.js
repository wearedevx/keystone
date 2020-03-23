const { flags } = require('@oclif/command')
const chalk = require('chalk')

const {
  listAllFiles,
  listEnvFiles,
} = require('@keystone.sh/core/lib/commands/list')

const { CommandSignedIn } = require('../lib/commands')

class ListCommand extends CommandSignedIn {
  async run() {
    const { flags } = this.parse(ListCommand)

    await this.withUserSession(async userSession => {
      const env = await this.getProjectEnv()
      const project = await this.getProjectName()

      if (flags.all) {
        await listAllFiles(userSession)
      } else {
        await listEnvFiles(userSession, { project, env })
      }
    })
  }
}

ListCommand.description = `list files tracked for your current environment
`

ListCommand.examples = [chalk.blue('$ ks list')]

ListCommand.flags = {
  ...CommandSignedIn.flags,
  all: flags.boolean({
    char: 'a',
    multiple: false,
    description: 'list every files in your gaia hub',
  }),
}

module.exports = ListCommand
