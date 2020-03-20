const chalk = require('chalk')
const {
  assertUserIsAdminOrContributor,
} = require('@keystone.sh/core/lib/member')
const { removeEnvFiles } = require('@keystone.sh/core/lib/env')
const { removeEnvFromProject } = require('@keystone.sh/core/lib/projects')

const { CommandSignedIn } = require('../../lib/commands')

class EnvCommand extends CommandSignedIn {
  /**
   * Remove env.
   * @param {*} project
   * @param {*} name
   */
  async removeEnv(project, name) {
    await this.withUserSession(async userSession => {
      await assertUserIsAdminOrContributor(userSession, { project })

      const absoluteProjectPath = await this.getConfigFolderPath()

      await removeEnvFiles(userSession, {
        project,
        env: name,
        absoluteProjectPath,
      })
      // Remove en from project descriptor.
      await removeEnvFromProject(userSession, {
        project,
        env: name,
      })
    })
  }

  async run() {
    const { args } = this.parse(EnvCommand)
    const project = await this.getProjectName()

    try {
      await this.removeEnv(project, args.env)
    } catch (error) {
      this.log(`${chalk.red(error)}`)
    }
  }
}

EnvCommand.description = `remove an environment`

EnvCommand.args = [
  {
    name: 'env',
    required: true, // make the arg required with `required: true`
    description: 'environment name', // help description
    hidden: false,
  },
]

EnvCommand.examples = [chalk.blue(`$ ks env rm ${chalk.italic('ENV_NAME')}`)]

module.exports = EnvCommand
