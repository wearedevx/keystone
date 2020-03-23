const chalk = require('chalk')
const { cli } = require('cli-ux')
const { deleteFiles } = require('@keystone.sh/core/lib/commands/delete')

const { CommandSignedIn } = require('../lib/commands')

class RmCommand extends CommandSignedIn {
  async deleteFile(project, env, files) {
    await this.withUserSession(async userSession => {
      cli.action.start('Deleting')
      let success = true
      const absoluteProjectPath = await this.getConfigFolderPath()
      try {
        const fileRelativePaths = await Promise.all(
          files.map(e => this.getFileRelativePath(e))
        )
        await deleteFiles(userSession, {
          project,
          env,
          files: fileRelativePaths,
          absoluteProjectPath,
        })
      } catch (error) {
        console.error(error)
        this.log(`▻ ${chalk.red(error.message)}\n`)
        success = false
      }
      cli.action.stop(success ? 'success' : 'failure')
      if (success) {
        files.map(file =>
          this.log(`> ${file} successfully deleted ${chalk.green.bold('✔')}`)
        )
      }
    })
  }

  async run() {
    const { argv } = this.parse(RmCommand)
    try {
      const project = await this.getProjectName()
      const env = await this.getProjectEnv()
      await this.deleteFile(project, env, argv)
    } catch (error) {
      this.log(error.message)
    }
  }
}

RmCommand.description = `Deletes one or more files.
If you're an administrator or a contributor, the files will be removed for everyone.
If you're a reader on the environment, you can't delete any files.
`

RmCommand.strict = false

RmCommand.args = [
  {
    name: 'filepaths',
    required: true, // make the arg required with `required: true`
    description: 'Path to your file. Accepts a glob pattern', // help description
    hidden: false,
  },
]

RmCommand.examples = [chalk.blue('$ ks rm path/to/file')]

module.exports = RmCommand
