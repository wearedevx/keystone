const { cli } = require('cli-ux')
const chalk = require('chalk')
const { flags } = require('@oclif/command')

const pull = require('@keystone/core/lib/commands/pull')

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

        pulledFiles.map(
          async ({ fileDescriptor, updated, descriptorUpToDate }) => {
            if (descriptorUpToDate) {
              this.log(`▻ You are already up to date. Nothing to do !`)
              return
            }
            if (updated) {
              this.log(
                `▻ File written to ${fileDescriptor.name} ${chalk.green.bold(
                  '✓'
                )}`
              )
            } else {
              this.log(
                `▻ File ${
                  fileDescriptor.name
                } already is the latest version ${chalk.green.bold('✓')}`
              )
            }
            // this.log(
            //   `▻ Couldn't save the file : ${chalk.bold(
            //     fileDescriptor.content.name
            //   )} ${chalk.red.bold('✗')}`
            // )
          }
        )
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
