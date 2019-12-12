const chalk = require('chalk')
const { cli } = require('cli-ux')
const { flags } = require('@oclif/command')
const util = require('util')

const { CommandSignedIn } = require('../lib/commands')
const {
  getFiles,
  getFileDescriptor,
  getFileFromGaia,
} = require('../lib/core/file')

class CatCommand extends CommandSignedIn {
  
  async cat(files, project, env, flags) {
    await this.withUserSession(async userSession => {
      cli.action.start('Fetching')
      let success
      try {
        let fetchedFiles
        if (flags.debug) {
          fetchedFiles = JSON.parse(await getFileFromGaia(userSession, files, { username: flags.origin }))
          console.log(util.inspect(fetchedFiles, false, null))
          return
        }

        const userData = userSession.loadUserData()

        const envDescriptor = await getFileDescriptor(userSession, {
          blockstack_id: userData.username,
          author: userData.username,
          type: 'env',
          project,
          env,
        })

        fetchedFiles = await getFiles(userSession, {
          project,
          files: [files],
          envDescriptor,
        })

        fetchedFiles.map(file =>
          file.fetched
            ? console.log(`${file.descriptor.content}\n`)
            : console.log(file)
        )
        success = true
      } catch (err) {
        console.log(err)
        success = err
      }
      cli.action.stop(success ? 'done' : success)
    })
  }

  async run() {
    try {
      const env = await this.getProjectEnv()
      const { args, flags } = this.parse(CatCommand)

      // at least 1 arguments required, an email
      // const project = await this.getProjectName(flags)
      const project = await this.getProjectName(flags)
      await this.cat(args.path, project, env, flags)
    } catch (error) {
      this.log(error.message)
    }
  }
}

CatCommand.args = [
  {
    name: 'path',
    required: true, // make the arg required with `required: true`
    description: 'path to your file', // help description
    hidden: false,
  },
  {
    name: 'decrypt',
    required: false, // make the arg required with `required: true`
    description: 'should decrypt', // help description
    default: true,
    hidden: false,
  },
]
CatCommand.flags = {
  ...CommandSignedIn.flags,
  debug: flags.boolean({
    char: 'd',
    multiple: false,
    default: false,
    description: `cat file with full path`,
  }),
  origin: flags.string({
    char: 'o',
    multiple: false,
    default: null,
    description: `from origin`,
  }),
  removal: flags.boolean({
    multiple: false,
    default: false,
    description: `Deletes an invitation`,
  }),
}

CatCommand.description = `Output a remote file.
`

CatCommand.examples = [chalk.yellow('$ ks cat my-file ')]

module.exports = CatCommand
