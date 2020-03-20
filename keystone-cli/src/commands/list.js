const { flags } = require('@oclif/command')
// const chalk = require('chalk')
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

ListCommand.description = `Lists projects, environments, members and files
`

ListCommand.examples = [chalk.blue('$ ks list files')]

ListCommand.flags = {
  ...CommandSignedIn.flags,
  all: flags.boolean({
    char: 'a',
    multiple: false,
    description:
      'For files listing, list every files in your gaia hub. For members, list files from project, instead of the environment.',
  }),
}

module.exports = ListCommand
