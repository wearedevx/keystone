const chalk = require('chalk')

const { CommandSignedIn } = require('../lib/commands')
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
      modifiedFiles
        .filter(f => f.status !== 'ok')
        .map(file => console.log(file.path, ':', chalk.bold(file.status)))
    })
  }
}

StatusCommand.description = `Shows the status of tracked files
`

StatusCommand.examples = [chalk.yellow('$ ks status')]

module.exports = StatusCommand
