const { flags } = require('@oclif/command')
const chalk = require('chalk')

const { CommandSignedIn } = require('../lib/commands')

const {
  removeFromEnv,
  removeFromProject,
} = require('@keystone.sh/core/lib/commands/remove')

class RemoveCommand extends CommandSignedIn {
  async removeUser({ project, env, users }) {
    await this.withUserSession(async userSession => {
      try {
        await Promise.all(
          users.map(user => {
            if (project) {
              return removeFromProject(userSession, {
                project,
                user,
              })
            } else if (env) {
              return removeFromEnv(userSession, {
                project,
                env,
                user,
              })
            } else {
              this.log('/!\\ Use -e to specify an environement.')
            }
          })
        )
      } catch (error) {
        console.error(error)
        this.log(`▻ ${chalk.red(error.message)}\n`)
      }
    })
  }

  async run() {
    let {
      flags: { project, env, users },
    } = this.parse(RemoveCommand)

    console.log(project, env, users)
    if (project) project = await this.getProjectName()
    try {
      await this.removeUser({ project, env, users })
    } catch (error) {
      this.log(error.message)
    }
  }
}

RemoveCommand.description = `Remove a user.
...
If you are an administrator, you can remove a user from a project.
`

RemoveCommand.flags = {
  ...CommandSignedIn.flags,
  env: flags.string({
    char: 'e',
    multiple: false,
    description: 'Environment you want to remove the user from.',
  }),

  project: flags.boolean({
    char: 'p',
    multiple: false,
    description: 'Remove the user from the project.',
  }),

  users: flags.string({
    char: 'u',
    multiple: true,
    description: 'List of user you want to remove. Separated by space.',
  }),
}

module.exports = RemoveCommand
