const { cli } = require('cli-ux')
const chalk = require('chalk')
const { checkoutEnv } = require('@keystone.sh/core/lib/env')

const { CommandSignedIn, execPull } = require('../../lib/commands')

class EnvCommand extends CommandSignedIn {
  async checkout(project, env) {
    await this.withUserSession(async userSession => {
      try {
        const absoluteProjectPath = await this.getConfigFolderPath()

        try {
          cli.action.start('Changing environment')
          await checkoutEnv(userSession, {
            project,
            env,
            absoluteProjectPath,
          })
          cli.action.stop('done')
        } catch (err) {
          if (err.code === 'PendingModification') {
            cli.action.stop('aborted')
            console.log(
              'You have modified files is your working directory. Please push your changes or use ',
              chalk.blue('$ ks env reset'),
              'to abort your changes.'
            )
            console.log('\n')
            err.data.forEach(f => console.log(f.path, '', chalk.bold(f.status)))

            process.exit(0)
          }
        }

        cli.action.stop('done')

        cli.action.start('Fetching files')
        await execPull(userSession, {
          project,
          env,
          absoluteProjectPath,
          force: true,
        })
        cli.action.stop('done')
      } catch (err) {
        cli.action.stop('failed')
        this.log(err)
      }
    })
  }

  async run() {
    const { args } = this.parse(EnvCommand)
    const project = await this.getProjectName()

    try {
      await this.checkout(project, args.env)
    } catch (error) {
      this.log(`${chalk.red(error)}`)
    }
  }
}

EnvCommand.description = `switch environment and pull files`

EnvCommand.args = [
  {
    name: 'env',
    required: true, // make the arg required with `required: true`
    description: 'environment name', // help description
    hidden: false,
  },
]

EnvCommand.examples = [
  chalk.blue(`$ ks env checkout ${chalk.italic('ENV_NAME')}`),
]

module.exports = EnvCommand
