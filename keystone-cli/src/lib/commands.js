const { getFile } = require('blockstack/lib/storage')
const { lookupProfile } = require('blockstack/lib/profiles/profileLookup')
const axios = require('axios')
const { decryptECIES } = require('blockstack/lib/encryption/ec')

const { Command, flags } = require('@oclif/command')
const { cli } = require('cli-ux')
const chalk = require('chalk')
const normalize = require('normalize-path')
const treeify = require('treeify')
const inquirer = require('inquirer')
const path = require('path')
const pull = require('@keystone.sh/core/lib/commands/pull')
const { pullShared } = require('@keystone.sh/core/lib/commands/share')
const { findProjectByUUID } = require('@keystone.sh/core/lib/projects')
const {
  isOneOrMoreAdmin,
  setMembersToEnvs,
} = require('@keystone.sh/core/lib/env')
const { logo } = require('../lib/ux')
const { KEYSTONE_WEB } = require('@keystone.sh/core/lib/constants')
const { getSession, getProjectConfig } = require('../lib/blockstackLoader')

const createSharedUserSession = token => {
  return {
    loadUserData: () => ({}),
    getFile,
  }
}

function promptEnvToChange(envs, project = false) {
  return inquirer.prompt([
    {
      type: 'list',
      name: 'env',
      message: `Which ${
        project ? 'project' : 'environment'
      } you want to configure?`,
      choices: [
        ...envs,
        new inquirer.Separator(),
        chalk.green('Save & quit'),
        chalk.red('Cancel'),
      ],
    },
  ])
}

function promptUsersToChange(role, env, allMembers, envsMembers, envs) {
  return inquirer.prompt([
    {
      type: 'checkbox',
      name: 'members',
      message: `Which users you want to make ${chalk.green(
        role
      )} for the ${chalk.green(env)} environment?`,
      choices: [
        ...allMembers.map(user => {
          if (
            envsMembers[env][role].find(
              member => member.blockstack_id === user.blockstack_id
            )
          ) {
            return { name: user.blockstack_id, value: user, checked: true }
          }
          return { name: user.blockstack_id, value: user }
        }),
      ],
      default: [envs],
    },
  ])
}

function promptRoleToChange(envsMembers, env) {
  return inquirer.prompt([
    {
      type: 'list',
      name: 'role',
      message: 'Which role you want to set?',
      choices: [
        {
          value: 'readers',
          name: `readers ${chalk.gray(`(${envsMembers[env].readers.length})`)}`,
        },
        {
          value: 'contributors',
          name: `contributors ${chalk.gray(
            `(${envsMembers[env].contributors.length})`
          )}`,
        },
        {
          value: 'admins',
          name: `admins ${chalk.gray(`(${envsMembers[env].admins.length})`)}`,
        },
        new inquirer.Separator(),
        chalk.red('Cancel'),
      ],
    },
  ])
}

const CommandSignedIn = class extends Command {
  async configureMembers(props) {
    try {
      let { envsMembers, currentStep, env, role } = props
      const { allMembers, project } = props

      const envs = Object.keys(envsMembers)

      // Nicely log env config with color
      let members = JSON.parse(JSON.stringify(envsMembers))

      members = Object.keys(members).reduce((acc, role) => {
        acc[chalk.red(role)] = Object.keys(members[role]).reduce(
          (acc2, curr2) => {
            acc2[chalk.blue(curr2)] = members[role][curr2].map(x =>
              chalk.green(x.blockstack_id)
            )
            return acc2
          },
          {}
        )
        return acc
      }, {})

      console.log(treeify.asTree(members, true))
      if (!isOneOrMoreAdmin(envsMembers)) {
        this.log(
          chalk.red(
            '/!\\ There must be one or more admin in each environment\n'
          )
        )
      }

      if (currentStep === 0) {
        const result = await promptEnvToChange(envs, project)

        if (/.*cancel.*/i.test(result.env)) {
          if (
            await CommandSignedIn.confirm(
              chalk.red('Sure you want to quit and abort all changes ?')
            )
          ) {
            process.exit(0)
          }
        } else if (/.*save.*/i.test(result.env)) {
          if (
            await CommandSignedIn.confirm(
              chalk.green('Sure you want to save your changes and quit ?')
            )
          ) {
            if (isOneOrMoreAdmin(envsMembers)) {
              return envsMembers
            }
          }
        } else {
          env = result.env
          currentStep = 1
        }
      }
      if (currentStep === 1) {
        const result = await promptRoleToChange(envsMembers, env)
        role = result.role

        if (/.*ancel.*/i.test(role)) {
          currentStep = 0
        } else {
          currentStep = 2
        }
      }
      if (currentStep === 2) {
        let { members } = await promptUsersToChange(
          role,
          env,
          allMembers,
          envsMembers,
          envs
        )

        currentStep = 0

        // members = members.map(member => ({
        //   blockstack_id: member,
        // }))
        envsMembers = setMembersToEnvs({ envsMembers, members, role, env })
      }
      console.log('\x1Bc')

      return this.configureMembers({
        allMembers,
        envsMembers,
        envs,
        currentStep,
        env,
        role,
        project,
      })
    } catch (err) {
      console.error(err)
    }
  }

  // async fetchProject(userSession, project) {
  //   try {
  //     cli.action.start('Fetching project')
  //     return await getProjectDescriptor(userSession, { project })
  //   } catch (error) {
  //     cli.action.stop('Failed')
  //     this.log(`▻ ${error.message} ${chalk.red.bold('✗')}`)
  //     // end the cli
  //     return false
  //   }
  // }

  async getConfigFolderPath() {
    const { absoluteProjectPath } = await getProjectConfig()
    return absoluteProjectPath
  }

  async getProjectEnv() {
    try {
      const {
        config: { env },
      } = await getProjectConfig()
      return env
    } catch (error) {
      this.log(
        `▻ Keystone config file is missing or malformed, please start with ${chalk.yellow(
          `$ ks init`
        )}`
      )
      if (!process.env.KEYSTONE_SHARED) process.exit(1)
    }
  }

  async getProjectName() {
    try {
      const {
        config: { project },
      } = await getProjectConfig()

      return project
    } catch (error) {
      console.log(
        `▻ Keystone config file is missing or malformed, please start with ${chalk.yellow(
          `$ ks init`
        )}`
      )
      if (!process.env.KEYSTONE_SHARED) {
        process.exit(1)
      }
    }
  }

  async getDefaultDirectory(flags) {
    if (flags && flags.directory) {
      return flags.directory
    }
    try {
      const { config } = await getProjectConfig()
      return config.directory
    } catch (error) {
      this.log(
        `▻ Keystone config file is missing or malformed, please start with ${chalk.yellow(
          `$ ks init`
        )}`
      )
      process.exit(0)
    }
  }

  async withUserSession(callback, options = {}) {
    const { configPath, config } = options
    const userSession = await this.getUserSession(configPath)
    if (userSession) {
      await callback(userSession)
    }
  }

  async getUserSession(configPath) {
    try {
      if (process.env.KEYSTONE_SHARED) {
        return createSharedUserSession(process.env.KEYSTONE_SHARED)
      }
      const config = configPath || this.config.configDir
      const userSession = await getSession(config)
      if (userSession && userSession.isUserSignedIn()) {
        return userSession
      }
      this.log(
        `▻ You're not connected, please sign in first: ${chalk.yellow(
          `$ ks login`
        )}`
      )
    } catch (error) {
      this.log(
        `▻ You're not connected, please sign in first: ${chalk.yellow(
          `$ ks login`
        )}`
      )
    }
    return false
  }

  async init() {
    this.log(logo)
  }

  async getFileRelativePath(filePath) {
    const absoluteFilePath = path.resolve(filePath)
    const absoluteConfigFolderPath = await this.getConfigFolderPath()

    if (absoluteFilePath.indexOf(absoluteConfigFolderPath) === -1) {
      throw new Error(
        `${filePath} is not in the keystone project ${absoluteConfigFolderPath}`
      )
    }

    const relativeFilePath = absoluteFilePath.replace(
      `${absoluteConfigFolderPath}/`,
      ''
    )

    return relativeFilePath
  }

  async getAbsolutePath(filePath) {
    const absoluteConfigFolderPath = await this.getConfigFolderPath()
    const absoluteFilePath = path.join(absoluteConfigFolderPath, filePath)
    return absoluteFilePath
  }
}

CommandSignedIn.normalizePath = (directory, filename) => {
  let dir = directory
  if (!directory) {
    dir = './'
  }
  return normalize(`${dir}/${filename}`)
}

CommandSignedIn.confirm = async message => {
  const answer = await inquirer.prompt([
    {
      name: 'confirm',
      type: 'confirm',
      message,
    },
  ])
  return answer.confirm
}

CommandSignedIn.flags = {
  project: flags.string({
    char: 'p',
    multiple: false,
    description: 'Set the project',
  }),
  config: flags.string({
    char: 'c',
    multiple: false,
    description: 'Set the path to the blockstack session file',
  }),
}

const getFileFromOne = async ({ privateKey, blockstackId, opts, filename }) => {
  const profile = await lookupProfile(opts.username)
  const apphub = profile.apps[KEYSTONE_WEB]
  const uri = `${apphub}${filename}`
  let file
  try {
    file = await axios.get(uri)
  } catch (error) {
    // we ignore 404 errors
    return false
  }

  const keyfileUnencrypted = await decryptECIES(privateKey, file.data)

  return keyfileUnencrypted
}

const execPull = async (
  userSession,
  { project, env, absoluteProjectPath, force }
) => {
  try {
    if (process.env.KEYSTONE_SHARED) {
      const buff = new Buffer(process.env.KEYSTONE_SHARED, 'base64')
      const { env, project, members, privateKey } = JSON.parse(
        buff.toString('ascii')
      )
      userSession = {
        ...userSession,
        getFile: (filename, opts) =>
          getFileFromOne({
            privateKey,
            blockstackId: members[0].blockstack_id,
            filename,
            opts,
          }),
        sharedPrivateKey: privateKey,
      }
      userSession.sharedPrivateKey = privateKey

      return pullShared(userSession, {
        absoluteProjectPath,
        project,
        env,
        origins: members,
      })
    }

    const pulledFiles = await pull(userSession, {
      project,
      env,
      absoluteProjectPath,
      force,
    })
    pulledFiles.map(
      async ({ fileDescriptor, updated, descriptorUpToDate, conflict }) => {
        if (descriptorUpToDate) {
          console.log(`▻ You are already up to date. Nothing to do !`)
          console.log(
            `▻ If you want to override your local files, use ${chalk.bold(
              --force
            )} flag`
          )
          return
        }
        if (updated) {
          if (!(typeof conflict === 'boolean')) {
            console.log(
              ` ${chalk.green.bold('✔')} ${fileDescriptor.name}: updated.`
            )
          } else if (conflict) {
            console.log(
              ` ${chalk.red.bold('✗')} ${
                fileDescriptor.name
              }: conflict. Correct them and push your changes !`
            )
          } else {
            console.log(
              ` ${chalk.green.bold('✔')} ${fileDescriptor.name}: auto-merge.`
            )
          }
        }
      }
    )
  } catch (error) {
    switch (error.code) {
      case 'PullWhileFilesModified':
        console.log(
          `Your files are modified. Please push your changes or re-run this command with --force to overwrite.`
        )
        error.data.forEach(file =>
          console.log(
            `▻ ${chalk.bold(file.path)} - ${file.status} ${chalk.red.bold('✗')}`
          )
        )
        break

      case 'MissingEnv':
        console.log(
          `You have no environment informed in your config.\nRun : ${chalk.blue(
            `ks env checkout ${chalk.italic('env_name')}`
          )}\n\nAvailable options :`
        )
        error.data.envs.forEach(env => console.log(`▻ ${chalk.bold(env)}`))
        break
      case 'FailedToFetch':
        console.log(
          `${chalk.red('Failed to fetch')} ${chalk.bold(
            "You don't have access to the environment."
          )}\n\nPlease ask your project administrator to give you access.\nRun : ${chalk.blue(
            `ks env checkout ${chalk.italic('env_name')}`
          )} to fetch files from another environment.`
        )
        break

      default:
        throw error
    }
  }
}

module.exports = {
  CommandSignedIn,
  execPull,
}
