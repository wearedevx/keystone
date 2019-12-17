const chalk = require('chalk')
const fs = require('fs')
const { newShare } = require('@keystone/core/lib/commands/share')

const { CommandSignedIn } = require('../lib/commands')

class ShareCommand extends CommandSignedIn {
  async newShare() {
    await this.withUserSession(async userSession => {
      const project = await this.getProjectName()
      const addedShare = newShare(userSession, { project })
    })
  }

  async run() {
    const { args } = this.parse(ShareCommand)

    if (args.action === 'new') {
      this.newShare()
    } else {
      this.pull()
    }
  }
}

ShareCommand.description = `Share your file file with a non-blockstack user
`

ShareCommand.examples = [chalk.yellow('$ ks share')]

ShareCommand.args = [
  {
    name: 'action',
    required: false, // make the arg required with `required: true`
    description: 'new || Path to config file', // help description
    hidden: false,
  },
]

module.exports = ShareCommand
