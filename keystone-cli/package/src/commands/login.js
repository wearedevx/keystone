const {Command, flags} = require('@oclif/command')
const { cli } = require('cli-ux')

const open = require('open')
const inquirer = require('inquirer')
const EC = require('elliptic').ec
const ec = new EC('secp256k1')
const blockstack = require('blockstack')
const { decryptECIES } = require('blockstack/lib/encryption/ec')
const axios = require('axios')
const {createUserSession} = require('../lib/blockstackLoader')
const { write } = require('../lib/cliStorage')



class LoginCommand extends Command {
  // static args = [{name: 'file'}]
  // static args = [{
  //     name: 'blockstack_id',
  //     required: false,            // make the arg required with `required: true`
  //     description: 'Your blockstack id', // help description
  //     hidden: false,
  //   }]
  

  async openLink(id) {
    const keypair = ec.genKeyPair()
    const pubPoint = keypair.getPublic()
    const publicKey = pubPoint.encode('hex');
    const _this = this    
    await open(`http://localhost:8000/confirm?token=${publicKey}&id=${id}`);
    cli.action.start('Linking your blockstack account...')
    const interval = setInterval(async () => {
      const keyfile = await _this.connect(id, publicKey)
      if(keyfile){
        clearInterval(interval)
        // Blockstack use the private in HEX to decrypt, see below
        // https://github.com/blockstack/blockstack.js/blob/master/src/encryption/ec.ts
        const keyfileUnencrypted = decryptECIES(keypair.getPrivate('hex'), keyfile.data)
        const userCredentials = JSON.parse(keyfileUnencrypted)
        const res = await write({
          path: this.config.configDir,
          filename: "session.json",
          content: keyfileUnencrypted
        })
        const userSession = createUserSession(userCredentials)
      }
    }, 3000)
    
    this.log(`login ${id} from ./src/commands/login.js`)
  }

  async connect(id, publicKey){
    const profile = await blockstack.lookupProfile(id)
    if(profile && profile.apps){
      const apphub = profile.apps['http://localhost:8000']
      if(apphub){
        try {
          const uri = `${apphub}${publicKey}.json`
          console.log("uri?", uri)
          const keyfile = await axios.get(uri)
          console.log("get keyfile")
          return keyfile
        } catch (error) {
          console.log("get keyfile -error", error.message)
          return false
        }
      }
      return false
    }
  }

  async prompt(){
    const answer = await inquirer.prompt([{
      name: 'blockstack_id',
      message: 'Enter your blockstack id:'
    }])
    if(answer.blockstack_id.length > 1){
      this.openLink(answer.blockstack_id)
    }
    else {
      await this.prompt()
    }
  }

  async run() {
    console.log("this.args", this.args)
    const { args } = this.parse(LoginCommand)
    if(args.blockstack_id){
      openLink(args.blockstack_id)
    }else{
      await this.prompt()
    }
  }
}

LoginCommand.description = `Describe the command here
...
Extra documentation goes here
`


// LoginCommand.flags = {
//   name: flags.string({char: 'n', description: 'name to print'}),
// }

module.exports = LoginCommand
