const chalk = require('chalk')
const { cli } = require('cli-ux')
const { flags } = require('@oclif/command')
const util = require('util')
const { readFileFromGaia } = require('@keystone.sh/core/lib/file/gaia')

const { CommandSignedIn } = require('../lib/commands')

class CatCommand extends CommandSignedIn {
  async cat(path, project, env, flags) {
    await this.withUserSession(async userSession => {
      cli.action.start('Fetching')
      let success
      try {
        let fetchedFiles
        if (flags.debug) {
          const opts = {
            origin: flags.origin,
            path,
            decrypt: flags.decrypt,
            json: flags.json,
          }

          fetchedFiles = await readFileFromGaia(userSession, opts)

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
          files: [path],
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
      cli.action.stop(success ? 'done' : 'failed')
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
  decrypt: flags.boolean({
    default: true,
    description: `Indiciate to decrypt or not`,
    allowNo: true,
  }),
  json: flags.boolean({
    multiple: false,
    default: true,
    description: `Indiciate to parse json or not`,
    allowNo: true,
  }),
}

CatCommand.description = `Output a remote file.
`

CatCommand.examples = [chalk.yellow('$ ks cat my-file ')]

module.exports = CatCommand
