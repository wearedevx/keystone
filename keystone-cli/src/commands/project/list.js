const chalk = require('chalk')

const { listProjects } = require('@keystone.sh/core/lib/commands/list')

const { CommandSignedIn } = require('../../lib/commands')

class ListCommand extends CommandSignedIn {
  async run() {
    await this.withUserSession(async userSession => {
      const projects = await listProjects(userSession)
      if (projects && projects.length > 0) {
        projects.forEach(project => {
          this.log(
            `> ${project.name.split('/')[0]}/${chalk.grey(
              project.name.split('/')[1]
            )}`
          )
        })
      } else {
        this.log('No project Found in user workspace!')
      }
    })
  }
}

ListCommand.description = `List projects in user workspace `

ListCommand.examples = [chalk.blue('$ ks project list')]

module.exports = ListCommand
