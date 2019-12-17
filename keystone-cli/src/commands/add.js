const { flags } = require('@oclif/command')
const chalk = require('chalk')
const { ROLES } = require('@keystone/core/lib/constants')
const { add } = require('@keystone/core/lib/commands/add')
const { assertUserIsAdmin } = require('@keystone/core/lib/member')
const { CommandSignedIn } = require('../lib/commands')

class AddCommand extends CommandSignedIn {
  async add(blockstackId, email, project) {
    const userSession = await this.getUserSession()

    await assertUserIsAdmin(userSession, { project })

    const invitee = {
      blockstackId,
      email,
      role: ROLES.READERS,
    }

    try {
      const memberAdded = await add(userSession, {
        invitee,
        project,
      })

      if (memberAdded.added) {
        this.log(
          `▻ ${chalk.yellow(
            memberAdded.blockstackId
          )} added to ${project} ${chalk.green.bold('✓')}`
        )
      } else {
        this.log(
          `▻ Failed to add ${chalk.yellow(
            memberAdded.blockstackId
          )} to ${project} ${chalk.red.bold('✗')}`,
          memberAdded.error
        )
      }
    } catch (error) {
      console.error(error)
      this.log(`\n${chalk.red(error.code)} : ${error.message}`)
    }
  }

  async run() {
    try {
      const { argv, flags, args } = this.parse(AddCommand)
      // at least 1 arguments required, a blockstack id
      if (argv.length >= 1) {
        const project = await this.getProjectName(flags)
        await this.add(args.blockstackId, args.email, project)
      } else {
        // await this.prompt()
        this.log(
          `We need at  least a blockstack id and an email associated to an invitation`
        )
      }
    } catch (error) {
      console.error(error)
      this.log(error.message)
    }
  }
}

AddCommand.args = [
  {
    name: 'blockstackId',
    required: true, // make the arg required with `required: true`
    description: 'Blockstack_id to add', // help description
    hidden: false,
  },
  {
    name: 'email',
    required: true, // make the arg required with `required: true`
    description: 'email associated to an invitation', // help description
    hidden: false,
  },
]

AddCommand.flags = {
  ...CommandSignedIn.flags,

  removal: flags.boolean({
    multiple: false,
    default: false,
    description: `Deletes an invitation`,
  }),
}

AddCommand.description = `Add a member to a project.

The member should have accepted your invitation for this to work
`

AddCommand.examples = [
  `${chalk.yellow(
    '$ ks add example.id.blockstack example@mail.com'
  )} ${chalk.gray.italic('#add a user to a project')}`,
]

module.exports = AddCommand
