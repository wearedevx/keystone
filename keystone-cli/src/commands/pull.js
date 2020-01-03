const { cli } = require('cli-ux')
const chalk = require('chalk')
const { flags } = require('@oclif/command')

const { CommandSignedIn, execPull } = require('../lib/commands')

class PullCommand extends CommandSignedIn {
  async pull({ project, env, force = false }) {
    await this.withUserSession(async userSession => {
      cli.action.start(`Fetching`)

      const absoluteProjectPath = await this.getConfigFolderPath()

      await execPull(userSession, {
        project,
        env,
        absoluteProjectPath,
        force,
      })
      cli.action.stop(`Done`)
    })
  }

  async run() {
    try {
      const { args, flags } = this.parse(PullCommand)
      const { force } = flags
      const project = await this.getProjectName()
      const env = await this.getProjectEnv()
      // const currentDirectory = await this.getDefaultDirectory()
      await this.pull({ project, env, force })
    } catch (error) {
      this.log(error)
    }
  }
}

PullCommand.flags = {
  force: flags.boolean({
    char: 'f',
    multiple: false,
    description: `Overwrite any changes`,
  }),
}

PullCommand.description = `Fetch files for current environment.`

module.exports = PullCommand
