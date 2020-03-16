const chalk = require('chalk')
const {
  getCacheFolder,
  getModifiedFilesFromCacheFolder,
} = require('@keystone.sh/core/lib/file/disk')

const { CommandSignedIn } = require('../lib/commands')

class StatusCommand extends CommandSignedIn {
  async run() {
    await this.withUserSession(async () => {
      const env = await this.getProjectEnv()
      const project = await this.getProjectName()
      const absoluteProjectPath = await this.getConfigFolderPath()
      const cacheFolder = getCacheFolder(absoluteProjectPath)

      const modifiedFiles = await getModifiedFilesFromCacheFolder(
        cacheFolder,
        absoluteProjectPath
      ).filter(f => f.status !== 'ok')

      if (env) console.log('On environment', chalk.bold(env))
      else console.log('No environment selected')

      console.log('Project', chalk.bold(project))
      console.log('\n')
      modifiedFiles.map(file =>
        console.log(
          '       ',
          chalk.red(chalk.bold(file.status), ':', file.path)
        )
      )
      if (modifiedFiles.length === 0) console.log('No file modified locally.')
      console.log('\n')
      console.log(
        `Use ${chalk.blue(
          `$ ks diff path/to/file`
        )} to see modification applied to a file.`
      )
    })
  }
}

StatusCommand.description = `Shows the status of tracked files
`

StatusCommand.examples = [chalk.blue('$ ks status')]

module.exports = StatusCommand
