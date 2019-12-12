const { cli } = require('cli-ux')
const chalk = require('chalk')
const { flags } = require('@oclif/command')

const pull = require('../lib/core-2.0/commands/pull')

const { CommandSignedIn } = require('../lib/commands')

class PullCommand extends CommandSignedIn {
  async pull({ project, env, force = false }) {
    await this.withUserSession(async userSession => {
      cli.action.start(`Fetching`)

      const absoluteProjectPath = await this.getConfigFolderPath()
      
      try {
        const pulledFiles = await pull(userSession, {
          project,
          env,
          absoluteProjectPath,
          force,
        })

        pulledFiles.map(async file => {
          if (file.name) {
            this.log(`▻ File written to ${file.name} ${chalk.green.bold('✓')}`)
          } else {
            this.log(
              `▻ Couldn't save the file : ${chalk.bold(file)} ${chalk.red.bold(
                '✗'
              )}`
            )
          }
        })
      } catch (error) {
        switch (error.code) {
          case 'PullWhileFilesModified':
            this.log(
              `Your files are modified. Please push your changes or re-run this command with --force to overwrite.`
            )
            error.data.forEach(file =>
              this.log(
                `▻ ${chalk.bold(file.path)} - ${file.status} ${chalk.red.bold(
                  '✗'
                )}`
              )
            )
            break
          default:
            throw error
        }
      }
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
