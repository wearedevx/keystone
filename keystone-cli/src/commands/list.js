const { flags } = require('@oclif/command')
const chalk = require('chalk')
const { cli } = require('cli-ux')

const { CommandSignedIn } = require('../lib/commands')
const {
  listEnvironments,
  listAllMembers,
  listEnvMembers,
  listAllFiles,
  listEnvFiles,
  listProjects,
} = require('@ks/core/lib/commands/list')

class ListCommand extends CommandSignedIn {
  async listProjects() {
    await this.withUserSession(async userSession => {
      cli.action.start('Fetching your projects')
      await listProjects(userSession)
      this.log(`\n`)
    })
    cli.action.stop('Done')
  }

  async run() {
    const { flags, args } = this.parse(ListCommand)

    await this.withUserSession(async userSession => {
      if (args.type === 'projects') {
        this.listProjects()
      } else {
        const env = await this.getProjectEnv()
        const project = await this.getProjectName()

        if (args.type === 'members') {
          if (flags.all) {
            listAllMembers(userSession, { project })
          } else {
            listEnvMembers(userSession, {
              project,
              env,
              isProjectMembers: false,
            })
          }
        } else if (args.type === 'files') {
          if (flags.all) {
            listAllFiles(userSession)
          } else {
            listEnvFiles(userSession, { project, env })
          }
        } else if (args.type === 'environments') {
          listEnvironments(userSession, { project })
        }
      }
    })
  }
}

ListCommand.description = `Lists projects, environments, members and files
`

ListCommand.examples = [
  chalk.yellow('$ ks list members'),
  // chalk.yellow('$ ks list --project=my-project'),
]

ListCommand.args = [
  {
    name: 'type',
    required: true, // make the arg required with `required: true`
    description:
      'What do you want to list (projects, environments, members or files)', // help description
    hidden: false,
  },
]

ListCommand.flags = {
  ...CommandSignedIn.flags,
  all: flags.boolean({
    char: 'a',
    multiple: false,
    description: 'List all elements',
  }),
}

module.exports = ListCommand
