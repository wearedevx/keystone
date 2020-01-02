const chalk = require('chalk')
const { cli } = require('cli-ux')
const deleteFiles = require('@ks/core/lib/commands/delete')

const { CommandSignedIn } = require('../lib/commands')

class DeleteCommand extends CommandSignedIn {
  async deleteFile(project, env, files) {
    await this.withUserSession(async userSession => {
      cli.action.start('Deleting')
      let success
      try {
        const fileRelativePaths = await Promise.all(
          files.map(e => this.getFileRelativePath(e))
        )

        await deleteFiles(userSession, {
          project,
          env,
          files: fileRelativePaths,
        })
      } catch (error) {
        console.error(error)
        this.log(`▻ ${chalk.red(error.message)}\n`)
        success = false
      }
      cli.action.stop(success ? 'success' : 'failure')
    })
  }

  async run() {
    const { argv } = this.parse(DeleteCommand)
    try {
      const project = await this.getProjectName()
      const env = await this.getProjectEnv()
      await this.deleteFile(project, env, argv)
    } catch (error) {
      this.log(error.message)
    }
  }
}

DeleteCommand.description = `Deletes one or more files.
If you're an administrator or a contributor, the files will be removed for everyone.
If you're a reader on the project, you can't delete any files.
`

DeleteCommand.strict = false

DeleteCommand.args = [
  {
    name: 'filepaths',
    required: false, // make the arg required with `required: true`
    description: 'Path to your file. Accepts a glob pattern', // help description
    hidden: false,
  },
]

module.exports = DeleteCommand
