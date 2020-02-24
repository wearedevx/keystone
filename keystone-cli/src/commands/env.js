const { cli } = require('cli-ux')
const chalk = require('chalk')
const { flags } = require('@oclif/command')
const inquirer = require('inquirer')
const {
  assertUserIsAdminOrContributor,
} = require('@keystone.sh/core/lib/member')
const {
  createEnv,
  removeEnvFiles,
  checkoutEnv,
} = require('@keystone.sh/core/lib/env')
const {
  getLatestMembersDescriptor,
  getMembers,
  getLatestProjectDescriptor,
} = require('@keystone.sh/core/lib/descriptor')
const {
  addEnvToProject,
  removeEnvFromProject,
} = require('@keystone.sh/core/lib/projects')

const { CommandSignedIn, execPull } = require('../lib/commands')
const { config } = require('@keystone.sh/core/lib/commands/env')
const { ROLES } = require('@keystone.sh/core/lib/constants')

class EnvCommand extends CommandSignedIn {
  async saveChanges(project, envsDescriptor, type) {
    await this.withUserSession(async userSession => {
      try {
        config(userSession, { project, descriptors: envsDescriptor, type })
      } catch (err) {
        cli.action.stop('failed')
        this.log(err)
      }
    })
  }

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

  /**
   * Remove env.
   * @param {*} project
   * @param {*} name
   */
  async removeEnv(project, name) {
    await this.withUserSession(async userSession => {
      await assertUserIsAdminOrContributor(userSession, { project })

      const absoluteProjectPath = await this.getConfigFolderPath()

      await removeEnvFiles(userSession, {
        project,
        env: name,
        absoluteProjectPath,
      })
      // Remove en from project descriptor.
      await removeEnvFromProject(userSession, {
        project,
        env: name,
      })
    })
  }

  async configureEnv(project) {
    await this.withUserSession(async userSession => {
      const { username } = userSession.loadUserData()
      try {
        // Check if env already exists.
        const projectDescriptor = await getLatestProjectDescriptor(
          userSession,
          {
            project,
          }
        )

        const envMembersDescriptors = await Promise.all(
          projectDescriptor.content.env.map(async env => {
            const descriptor = await getLatestMembersDescriptor(userSession, {
              project,
              env,
            })

            return { env, descriptor }
          })
        )

        let envsMembers = envMembersDescriptors.reduce(
          (envs, { env, descriptor }) => {
            envs[env] = descriptor.content
            return envs
          },
          {}
        )

        // only environments admins can change users permissions
        // so we keep only environments where the user is admin
        envsMembers = Object.keys(envsMembers).reduce((envs, env) => {
          const isAdmin = envsMembers[env][ROLES.ADMINS].find(
            member => member.blockstack_id === username
          )
          if (isAdmin) {
            envs[env] = envsMembers[env]
          }
          return envs
        }, {})

        const allMembers = await getMembers(userSession, { project })

        this.log('\x1Bc')

        await this.configureMembers({
          allMembers,
          envsMembers,
          currentStep: 0,
          type: 'env',
        })

        await this.saveChanges(project, envMembersDescriptors)
      } catch (err) {
        this.log(chalk.bold(err))
      }
    })
  }

  async checkout(project, env) {
    await this.withUserSession(async userSession => {
      try {
        const absoluteProjectPath = await this.getConfigFolderPath()
        try {
          cli.action.start('Changing environment')
          await checkoutEnv(userSession, {
            project,
            env,
            absoluteProjectPath,
          })
          cli.action.stop('done')
        } catch (err) {
          cli.action.stop('error')
          if (err.code === 'PendingModification') {
            console.log('\n')
            err.data.forEach(f => console.log(f.path, '', f.status))
            console.log('\n')

            const { force } = await inquirer.prompt([
              {
                name: 'force',
                type: 'confirm',
                message:
                  'You have pending modification in tracked files. Do you want to override them ?',
              },
            ])
            if (force) {
              cli.action.start('Changing environment')

              await checkoutEnv(userSession, {
                project,
                env,
                absoluteProjectPath,
                force,
              })
            } else {
              cli.action.stop('aborted')

              process.exit(0)
            }
          }
        }
        cli.action.stop('done')

        cli.action.start('Fetching files')
        await execPull(userSession, {
          project,
          env,
          absoluteProjectPath,
          force: true,
        })
        cli.action.stop('done')
      } catch (err) {
        cli.action.stop('failed')
        this.log(err)
      }
    })
  }

  async run() {
    const { args } = this.parse(EnvCommand)
    const project = await this.getProjectName()

    try {
      if (args.action) {
        if (args.action === 'config') {
          this.configureEnv(project)
        } else if (args.action === 'new') {
          if (args.env) {
            await this.newEnv(project, args.env)
          } else {
            throw new Error(`You need to give the name of the environment`)
          }
        } else if (args.action === 'remove') {
          if (args.env) {
            await this.removeEnv(project, args.env)
          } else {
            throw new Error(`You need to give the name of the environment`)
          }
        } else if (args.action === 'checkout') {
          if (args.env) {
            await this.checkout(project, args.env)
          } else {
            throw new Error(`You need to give the name of the environment`)
          }
        }
      } else {
        this.log(
          `▻ Current environment : ${chalk.bold(await this.getProjectEnv())}`
        )
      }
    } catch (error) {
      this.log(`${chalk.red(error)}`)
    }
  }
}

EnvCommand.description = `Manage environments.

You need to be administrator in the project in order to access the command.
`

EnvCommand.args = [
  {
    name: 'action',
    required: false, // make the arg required with `required: true`
    description: `  - config
    Change users role for each environment.

    You can change the role set by using the role flag. You have 3 choices:
      - reader: can only read files from the the environment and pull them locally
      - contributor: can read, write and add new files to the environement
      - admin: all the above plus ask people to join the project

  - new 
    Create a new environment

  - remove 
    Remove an environment
    `,
    hidden: false,
  },
  {
    name: 'env',
    required: false, // make the arg required with `required: true`
    description: 'Set working env', // help description
    hidden: false,
  },
]

EnvCommand.examples = [
  chalk.yellow('$ ks env config'),
  chalk.yellow(`$ ks env new ${chalk.italic('ENV_NAME')}`),
  chalk.yellow(`$ ks env remove ${chalk.italic('ENV_NAME')}`),
]

module.exports = EnvCommand
