const chalk = require('chalk')
const fs = require('fs')
const path = require('path')
const { CommandSignedIn } = require('../lib/commands')
const { KEYSTONE_HIDDEN_FOLDER } = require('@keystone.sh/core/lib/constants')

class StatusCommand extends CommandSignedIn {
  async run() {
    await this.withUserSession(async userSession => {
      const env = await this.getProjectEnv()
      const project = await this.getProjectName()
      const cachePath = path.join(
        await this.getConfigFolderPath(),
        KEYSTONE_HIDDEN_FOLDER,
        'cache'
      )
      const modifiedFiles = fs.readdirSync(cachePath)

      console.log('On environment', chalk.bold(env))
      modifiedFiles.map(file => console.log(chalk.bold(file), ' : modified'))
      console.log('\n')
    })
  }
}

StatusCommand.description = `Shows the status of tracked files
`

StatusCommand.examples = [chalk.yellow('$ ks status')]

module.exports = StatusCommand
