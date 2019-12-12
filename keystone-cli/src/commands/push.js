const { CommandSignedIn } = require('../lib/commands')

const { readFileFromDisk } = require('@keystone/core/lib/file')

const { push, pushModifiedFiles } = require('@keystone/core/lib/commands/push')

class PushCommand extends CommandSignedIn {
  async push(project, env, filenames) {
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
        await push(userSession, {
          project,
          env,
          files,
          absoluteProjectPath,
        })
      } else {
        // Push all modified files.
        await pushModifiedFiles(userSession, {
          project,
          env,
          absoluteProjectPath,
        })
      }
    })
  }

  async run() {
    const { argv } = this.parse(PushCommand)
    try {
      // at least 1 arguments required, a glob that returns 1 file and the project name
      // if (argv.length >= 1) {
      const project = await this.getProjectName()
      const env = await this.getProjectEnv()
      await this.push(project, env, argv)
      // } else {
      //   throw new Error('You need to specify at least one filename!')
      // }
    } catch (error) {
      console.log('error', error)
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

module.exports = PushCommand
