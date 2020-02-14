// const inquirer = require('inquirer')
// const { flags } = require('@oclif/command')
// const chalk = require('chalk')
const { CommandSignedIn } = require('../lib/commands')

class CheckCommand extends CommandSignedIn {
//   async askEmail() {
//     const answer = await inquirer.prompt([
//       {
//         name: 'email',
//         type: 'input',
//         message: `What's your email address?`,
//       },
//     ])
//     return answer.email
//   }

//   async check() {
//     const userSession = await this.getUserSession()

//     try {
//       // const invitations = await getInvitations(userSession)

//       const pubkey = await getPubkey(userSession, {
//         blockstack_id: 'l_abigael.id.blockstack',
//       })
//       this.log('pubkey', pubkey)
//     } catch (error) {
//       this.log(error)
//       this.log(`${error.message}`)
//     }
//   }

//   async run() {
//     try {
//       const { argv, flags } = this.parse(CheckCommand)
//       //at least 1 arguments required, an email
//       // const project = await this.getProjectName(flags)
//       await this.check()
//     } catch (error) {
//       this.log(error.message)
//     }
//   }
}

// // CheckCommand.args = [
// //   {
// //     name: 'emails',
// //     required: true,            // make the arg required with `required: true`
// //     description: 'Emails for invitations to be sent', // help description
// //     hidden: false,
// //   }
// // ]

// // CheckCommand.flags = {
// //   ...CommandSignedIn.flags,
// //   role: flags.string({
// //     char: 'r',
// //     multiple: false,
// //     options: ['reader', 'contributor', 'admin'],
// //     default: 'reader',
// //     description: `Assigns a role`
// //   }),
// //   removal: flags.boolean({
// //     multiple: false,
// //     default: false,
// //     description: `Deletes an invitation`
// //   })
// // }

// CheckCommand.description = `Check if a blockstack id has a Keystone application public key.
// `

// CheckCommand.examples = [`${chalk.yellow('$ ks check example.id.blockstack')}`]

module.exports = CheckCommand
