const chalk = require('chalk')
const fs = require('fs')
const { CommandSignedIn } = require('../lib/commands')

class WorkingProject extends CommandSignedIn {
  async run() {
    await this.withUserSession(async userSession => {
      const { project } = JSON.parse(fs.readFileSync('.ksconfig'))
      this.log(`▻ Working project :  ${chalk.bold(project)}`)
    })
  }
}

WorkingProject.description = `Print project name in current workspace
`

WorkingProject.examples = [chalk.yellow('$ ks wp')]

module.exports = WorkingProject
