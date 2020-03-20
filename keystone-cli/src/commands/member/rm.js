const chalk = require('chalk')

const { removeFromProject } = require('@keystone.sh/core/lib/commands/remove')
const { CommandSignedIn } = require('../../lib/commands')

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
    const { argv: users } = this.parse(RemoveCommand)
    const project = await this.getProjectName()
    try {
      await this.removeUser({ project, users })
    } catch (error) {
      this.log(error.message)
    }
  }
}

RemoveCommand.strict = false
RemoveCommand.description = `remove one or more users`

RemoveCommand.examples = [
  chalk.blue('$ ks remove nickname1.id.blockstack nickname2.id.blockstack'),
]

module.exports = RemoveCommand
