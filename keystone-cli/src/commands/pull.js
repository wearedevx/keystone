const { cli } = require('cli-ux')
const { flags } = require('@oclif/command')
const chalk = require('chalk')

const { CommandSignedIn, execPull } = require('../lib/commands')

/**
 * Pull files from the designated project and environment.
 */
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
      cli.action.stop(`done`)
    })
  }

  async run() {
    try {
      if (process.env.KEYSTONE_SHARED) {
        await this.pull({})
      } else {
        const { flags } = this.parse(PullCommand)
        const { force } = flags
        const project = await this.getProjectName()
        const env = await this.getProjectEnv()
        await this.pull({ project, env, force })
      }
    } catch (error) {
      this.log(error)
    }
  }
}

PullCommand.flags = {
  force: flags.boolean({
    char: 'f',
    multiple: false,
    description: `Overwrite any changes made locally`,
  }),
}

PullCommand.examples = [chalk.blue('$ ks pull')]

PullCommand.description = `Fetch files for current environment. Write them locally.

Once pulled files can be one of the three states :
  - updated : The file has been updated because someone else pushed a newer version
  - auto-merged : The file was modified and has been merged with someone else's changes
  - conflicted : The file has been modified and some lines are in conflict with someone else's changes. You should fix the conflicts and push your changes
`

module.exports = PullCommand
