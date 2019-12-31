const chalk = require('chalk')
const { CommandSignedIn } = require('../lib/commands')

const shouldDelete = filename =>
  filename !== 'public.key' && filename !== 'public.key.sig'

class ResetCommand extends CommandSignedIn {
  // remove all files except public.key
  async removeAllFiles() {
    await this.withUserSession(async userSession => {
      userSession.listFiles(file => {
        if (shouldDelete(file)) userSession.deleteFile(file)
        return true
      })
    })
  }

  async run() {
    this.removeAllFiles()
  }
}

ResetCommand.description = `Remove everything but your public.key file.
`

ResetCommand.examples = [chalk.yellow('$ ks reset')]

module.exports = ResetCommand
