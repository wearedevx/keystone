const debug = require('debug')('keystone:command:project')
const { cli } = require('cli-ux')
const chalk = require('chalk')
const {
  assertUserIsAdminOrContributor,
} = require('@keystone.sh/core/lib/member')
const {
  getLatestMembersDescriptor,
  getMembers,
} = require('@keystone.sh/core/lib/descriptor')

const { config } = require('@keystone.sh/core/lib/commands/project')
const { CommandSignedIn } = require('../lib/commands')

class ProjectCommand extends CommandSignedIn {
  async saveChanges(project, projectDescriptor) {
    await this.withUserSession(async userSession => {
      try {
        config(userSession, {
          project,
          descriptor: projectDescriptor,
          type: 'project',
        })
      } catch (err) {
        cli.action.stop('failed')
        this.log(err)
      }
    })
  }

  async configureProject(project) {
    await this.withUserSession(async userSession => {
      await assertUserIsAdminOrContributor(userSession, { project })

      try {
        debug('Get last project descriptor')

        this.log('\x1Bc')

        const allMembers = await getMembers(userSession, { project })

        const projectMembersDescriptor = await getLatestMembersDescriptor(
          userSession,
          {
            project,
            type: 'members',
          }
        )

        await this.configureMembers({
          allMembers,
          projectMembers: projectMembersDescriptor.content,
          envsMembers: { [project]: projectMembersDescriptor.content },
          project: true,
          env: project,
          currentStep: 1,
          type: 'project',
        })

        await this.saveChanges(project, projectMembersDescriptor)
      } catch (err) {
        this.log(chalk.bold(err))
      }
    })
  }

  async run() {
    const { args } = this.parse(ProjectCommand)
    const project = await this.getProjectName()

    try {
      if (args.action) {
        if (args.action === 'config') {
          this.configureProject(project)
        }
      }
    } catch (error) {
      this.log(`${chalk.red(error)}`)
    }
  }
}

ProjectCommand.description = `Manage users role in the project.

    You can change the role set by using the role flag. You have 3 choices:
- reader: can't do anything regarding the project itself.
- contributor: can add or remove environments.
- administrator: can add or remove environments, add and remove users, change users roles.
`

ProjectCommand.args = [
  {
    name: 'action',
    required: true,
    description: 'Configure project members',
  },
]

ProjectCommand.examples = [chalk.blue('$ ks project config')]

module.exports = ProjectCommand
