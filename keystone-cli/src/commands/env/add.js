const chalk = require('chalk')
const {
  assertUserIsAdminOrContributor,
} = require('@keystone.sh/core/lib/member')
const { createEnv } = require('@keystone.sh/core/lib/env')
const {
  getLatestProjectDescriptor,
} = require('@keystone.sh/core/lib/descriptor')
const { addEnvToProject } = require('@keystone.sh/core/lib/projects')

const { CommandSignedIn } = require('../../lib/commands')

class EnvCommand extends CommandSignedIn {
  /**
   * Create env
   */
  async newEnv(project, name) {
    await this.withUserSession(async userSession => {
      await assertUserIsAdminOrContributor(userSession, { project })
      try {
        // Check if env already exists.
        const projectDescriptor = await getLatestProjectDescriptor(
          userSession,
          {
            project,
          }
        )

        // If not, create it.
        if (projectDescriptor.content.env.includes(name)) {
          throw new Error(`Env ${name} already exists.`)
        }

        await createEnv(userSession, {
          env: name,
          projectDescriptor,
        })

        await addEnvToProject(userSession, {
          projectDescriptor,
          env: name,
        })
        this.log(`▻ Environment ${chalk.bold(name)} successfully created`)
      } catch (err) {
        this.log(err)
        this.log(`▻ Environment creation failed : ${chalk.bold(err)}`)
      }
    })
  }

  async run() {
    const { args } = this.parse(EnvCommand)
    const project = await this.getProjectName()

    try {
      if (args.env) {
        await this.newEnv(project, args.env)
      } else {
        throw new Error(`You need to give the name of the environment`)
      }
    } catch (error) {
      this.log(`${chalk.red(error)}`)
    }
  }
}

EnvCommand.description = `Manage environments.

You need to be administrator in the project in order to access the command.

You can change the role set by using the role flag. You have 3 choices:
- reader: can only read files from the the environment and pull them locally
- contributor: can read, write and add new files to the environement
- admin: all the above plus ask people to join the project
`

EnvCommand.args = [
  {
    name: 'env',
    required: true, // make the arg required with `required: true`
    description: 'Set working env', // help description
    hidden: false,
  },
]

EnvCommand.examples = [chalk.blue(`$ ks env add ${chalk.italic('ENV_NAME')}`)]

module.exports = EnvCommand
