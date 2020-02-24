const chalk = require('chalk')
const fs = require('fs')
const path = require('path')
const { CommandSignedIn } = require('../lib/commands')
const { KEYSTONE_HIDDEN_FOLDER } = require('@keystone.sh/core/lib/constants')
const {
  getCacheFolder,
  getModifiedFilesFromCacheFolder,
} = require('@keystone.sh/core/lib/file/disk')

class StatusCommand extends CommandSignedIn {
  async run() {
    await this.withUserSession(async userSession => {
      const env = await this.getProjectEnv()
      const project = await this.getProjectName()
      const absoluteProjectPath = await this.getConfigFolderPath()
      const cacheFolder = getCacheFolder(absoluteProjectPath)

      const modifiedFiles = await getModifiedFilesFromCacheFolder(
        cacheFolder,
        absoluteProjectPath
      )

      console.log('On environment', chalk.bold(env))
      console.log('Project', chalk.bold(project))
      console.log('\n')
      modifiedFiles.map(file =>
        console.log(file.path, ':', chalk.bold(file.status))
      )
      console.log('\n')
    })
  }
}

StatusCommand.description = `Shows the status of tracked files
`

StatusCommand.examples = [chalk.yellow('$ ks status')]

module.exports = StatusCommand
