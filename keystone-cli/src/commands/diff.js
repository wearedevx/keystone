const chalk = require('chalk')
const diff = require('@keystone.sh/core/lib/commands/diff')

const { CommandSignedIn } = require('../lib/commands')

class DiffCommand extends CommandSignedIn {
  async name(params) {}

  async run() {
    const { args } = this.parse(DiffCommand)
    await this.withUserSession(async userSession => {
      const absoluteProjectPath = await this.getConfigFolderPath()
      const filePath = await this.getFileRelativePath(args.filepath)
      try {
        const output = await diff(userSession, {
          absoluteProjectPath,
          filePath,
          file: args.filepath,
        })
        console.log(output)
      } catch (err) {
        console.log(err.message)
      }
    })
  }
}

DiffCommand.description = `Show a diff file
`

DiffCommand.examples = [chalk.yellow('$ ks diff', chalk.yellow('filepath'))]

DiffCommand.args = [
  {
    name: 'filepath',
    required: true, // make the arg required with `required: true`
    description: 'Path to your file.', // help description
    hidden: false,
  },
]

module.exports = DiffCommand
