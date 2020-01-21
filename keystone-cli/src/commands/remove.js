const { flags } = require('@oclif/command')
const chalk = require('chalk')

const { CommandSignedIn } = require('../lib/commands')

const { removeFromProject } = require('@keystone.sh/core/lib/commands/remove')

class RemoveCommand extends CommandSignedIn {
  async removeUser({ project, users }) {
    await this.withUserSession(async userSession => {
      try {
        await Promise.all(
          users.map(user => {
            return removeFromProject(userSession, {
              project,
              user,
            })
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
      flags: { users },
    } = this.parse(RemoveCommand)

    const project = await this.getProjectName()
    try {
      await this.removeUser({ project, users })
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

  users: flags.string({
    char: 'u',
    multiple: true,
    description: 'List of user you want to remove. Separated by space.',
  }),
}

module.exports = RemoveCommand
