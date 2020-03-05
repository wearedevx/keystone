const { flags } = require('@oclif/command')
const inquirer = require('inquirer')
const chalk = require('chalk')
const { cli } = require('cli-ux')
const {
  deleteFiles,
  deleteProject,
} = require('@keystone.sh/core/lib/commands/delete')

const { CommandSignedIn } = require('../lib/commands')

class DeleteCommand extends CommandSignedIn {
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
    const { argv, flags } = this.parse(DeleteCommand)
    try {
      if (flags.project) {
        await this.deleteProject(flags.project)
        return
      }
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
If you're a reader on the environment, you can't delete any files.
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

DeleteCommand.flags = {
  project: flags.string({
    char: 'p',
    multiple: false,
    description: `This is a debug command.\nUse this flag to completely delete all files of a project from your storage.`,
  }),
}

module.exports = DeleteCommand
