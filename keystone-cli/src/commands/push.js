const chalk = require('chalk')
const { cli } = require('cli-ux')
const { flags } = require('@oclif/command')
const { readFileFromDisk } = require('@keystone.sh/core/lib/file')
const { writeFileToGaia } = require('@keystone.sh/core/lib/file/gaia')

const {
  push,
  pushModifiedFiles,
} = require('@keystone.sh/core/lib/commands/push')
const { CommandSignedIn } = require('../lib/commands')

class PushCommand extends CommandSignedIn {
  async push(project, env, filenames, flags) {
    await this.withUserSession(async userSession => {
      const absoluteProjectPath = await this.getConfigFolderPath()

      // Push only files in args
      if (filenames.length > 0) {
        const files = await Promise.all(
          filenames.map(async f => {
            return {
              filename: await this.getFileRelativePath(f),
              fileContent: await readFileFromDisk(f),
            }
          })
        )

        if (flags.debug) {
          this.log('debug')
          await Promise.all(
            files.map(async file => {
              await writeFileToGaia(userSession, {
                path: flags.debug,
                content: file.fileContent,
                encrypt: flags.encrypt,
              })
            })
          )
        }
        // TODO make it work with jest
        cli.action.start('Pushing into private locker')
        const pushedFiles = await push(userSession, {
          project,
          env,
          files,
          absoluteProjectPath,
        })

        pushedFiles.forEach(f => {
          this.log(
            `▻ File ${chalk.bold(f.filename)} ${
              f.updated
                ? 'successfully pushed'
                : 'already is the latest version'
            } ${chalk.green.bold('✓')}`
          )
        })

        cli.action.stop('done')
      } else {
        cli.action.start('Pushing into private locker')
        // Push all modified files.
        const pushedFiles = await pushModifiedFiles(userSession, {
          project,
          env,
          absoluteProjectPath,
        })
        if (pushedFiles && pushedFiles.length > 0) {
          pushedFiles.forEach(f => {
            this.log(
              `▻ File ${chalk.bold(f.filename)} ${
                f.updated
                  ? 'successfully pushed'
                  : 'already is the latest version'
              } ${chalk.green.bold('✓')}`
            )
          })
        }
        cli.action.stop('done')
      }
    })
  }

  async run() {
    const { argv, flags } = this.parse(PushCommand)
    try {
      // at least 1 arguments required, a glob that returns 1 file and the project name
      // if (argv.length >= 1) {
      const project = await this.getProjectName()
      const env = await this.getProjectEnv()

      await this.push(project, env, argv, flags)
      // } else {
      //   throw new Error('You need to specify at least one filename!')
      // }
    } catch (error) {
      this.log('error', error)
    }
  }
}

PushCommand.description = `Push a file to a project.
`

PushCommand.strict = false

PushCommand.args = [
  {
    name: 'filepath',
    required: false, // make the arg required with `required: true`
    description: 'Path to your file. Accepts a glob pattern', // help description
    hidden: false,
  },
]

PushCommand.flags = {
  ...CommandSignedIn.flags,
  path: flags.string({
    char: 'p',
    multiple: false,
    description: '* DEBUG ONLY * push the file to the given path',
  }),
  encrypt: flags.string({
    char: 'e',
    multiple: false,
    description: '* DEBUG ONLY * encrypt the file with given blockstackid',
  }),
}

module.exports = PushCommand
