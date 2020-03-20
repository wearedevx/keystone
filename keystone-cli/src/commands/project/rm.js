const inquirer = require('inquirer')
const chalk = require('chalk')
const { deleteProject } = require('@keystone.sh/core/lib/commands/delete')

const { CommandSignedIn } = require('../../lib/commands')

class RemoveCommand extends CommandSignedIn {
  async deleteProject(project) {
    try {
      await this.withUserSession(async userSession => {
        const message = `${chalk.red(
          'Are you absolutly sure ?'
        )}\nThis action cannot be undone. If you are the only administrator in the project, nobody will be able to work on the project anymore.
        \nType ${chalk.yellow(project)} to delete the project:`

        const { input } = await inquirer.prompt([
          {
            name: 'input',
            message,
          },
        ])
        if (input === project) await deleteProject(userSession, { project })
        else console.log('No match found, nothing has been done.')
      })
    } catch (err) {
      console.error(err)
    }
  }

  async run() {
    const { args: project } = this.parse(RemoveCommand)
    try {
      await this.deleteProject(project)
    } catch (error) {
      this.log(error.message)
    }
  }
}

RemoveCommand.description = `remove project and its files from your storage`

RemoveCommand.args = [
  {
    name: 'project',
    required: true, // make the arg required with `required: true`
    description: 'project name (with uuid)', // help description
    hidden: false,
  },
]

module.exports = RemoveCommand
