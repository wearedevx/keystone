const { Command } = require('@oclif/command')
const { cli } = require('cli-ux')
const open = require('open')
const inquirer = require('inquirer')
const EC = require('elliptic').ec
const { decryptECIES } = require('blockstack/lib/encryption/ec')
const axios = require('axios')
const chalk = require('chalk')
const {
  KEYSTONE_WEB,
  LOGIN_KEY_PREFIX,
  SESSION_FILENAME,
} = require('@keystone/core/lib/constants')
const {
  createUserSession,
  getFilepath,
  getAppHub,
} = require('../lib/blockstackLoader')
const { write } = require('../lib/cliStorage')

const { logo } = require('../lib/ux')

const ec = new EC('secp256k1')
class LoginCommand extends Command {
  openLink(id) {
    return new Promise(async resolve => {
      const apphub = await getAppHub(id)

      if (apphub) {
        const { publicKey, privateKey } = LoginCommand.getKeypair()
        const uri = getFilepath({
          apphub,
          filename: `${LOGIN_KEY_PREFIX}${publicKey}.json`,
        })

        await open(`${KEYSTONE_WEB}/confirm?token=${publicKey}&id=${id}`)

        cli.action.start('Linking your account...')

        // FIXME: we mock setInterval/clearInterval with Jest during the tests/
        // the function becomes asynchronous and takes a fakeInterval
        // parameters.
        const interval = setInterval(async () => {
          const keyfile = await LoginCommand.connect(uri)
          if (keyfile) {
            clearInterval(interval)

            // Blockstack use the private in HEX to decrypt, see below
            // https://github.com/blockstack/blockstack.js/blob/master/src/encryption/ec.ts
            const keyfileUnencrypted = await decryptECIES(
              privateKey,
              keyfile.data
            )

            const userCredentials = JSON.parse(keyfileUnencrypted)
            await write({
              path: this.config.configDir,
              filename: SESSION_FILENAME,
              content: keyfileUnencrypted,
            })

            const userSession = createUserSession(userCredentials)
            if (userSession && userSession.isUserSignedIn()) {
              // well done we're connected
              const userData = userSession.loadUserData()
              this.log(
                `▻ You are connected under ${chalk.bold(userData.username)}`
              )
              this.log(`▻ You can logout with: ${chalk.yellow(`$ ks logout`)}`)
              // remove every files used to connect the terminal
              try {
                userSession.listFiles(async file => {
                  if (file.indexOf(LOGIN_KEY_PREFIX) > -1) {
                    userSession.deleteFile(file)
                  }
                  return true
                })
              } catch (error) {
                this.log(
                  "Can't remove your temporary keyfile from Gaïa. Please open a Github issue.",
                  error.message
                )
              } finally {
                resolve()
              }
            }
            cli.action.stop('Done')
          }
        }, 3000)
      } else {
        this.log(
          `▻ Unable to resolve your blockstack id. ${chalk.yellow(
            `Check your account`
          )}`
        )
        resolve()
      }
    })
  }

  async prompt() {
    const answer = await inquirer.prompt([
      {
        name: 'blockstack_id',
        message: 'Enter your blockstack id:',
      },
    ])
    if (answer.blockstack_id.length > 1) {
      await this.openLink(answer.blockstack_id)
    } else {
      await this.prompt()
    }
  }

  async run() {
    this.log(logo)
    const { args } = this.parse(LoginCommand)
    try {
      if (args.blockstack_id) {
        await this.openLink(args.blockstack_id)
      } else {
        await this.prompt()
      }
    } catch (error) {
      console.log('error', error)
    }
  }
}

LoginCommand.getKeypair = () => {
  const keypair = ec.genKeyPair()
  const pubPoint = keypair.getPublic()

  return {
    publicKey: pubPoint.encode('hex'),
    privateKey: keypair.getPrivate('hex'),
  }
}

LoginCommand.connect = async uri => {
  try {
    const keyfile = await axios.get(uri)
    return keyfile
  } catch (error) {
    // we ignore 404 errors
    return false
  }
}

LoginCommand.description = `Logs into your account with Blockstack or creates a new one
`
LoginCommand.args = [
  {
    name: 'blockstack_id',
    required: false, // make the arg required with `required: true`
    description: 'Your blockstack id', // help description
    hidden: false,
  },
]

// LoginCommand.flags = {
//   name: flags.string({char: 'n', description: 'name to print'}),
// }

module.exports = LoginCommand
