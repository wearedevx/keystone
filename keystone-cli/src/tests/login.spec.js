const EC = require('elliptic').ec

// Mock "open" module only for this test.
// If you want to always mock module for all tests,
// create a file in __mocks__ folder (at root),
// with file name with the module name, like (__mocks__/open.js):
// module.exports = jest.genMockFromModule('open')
jest.mock('open', () => {
  return jest.genMockFromModule('open')
})

const open = require('open')
const LoginCommand = require('../commands/login')
const {
  login,
  logout,
  runCommand,
  putFile,
  getSessionWithConfig,
} = require('./helpers')

describe('Login Command', () => {
  let result
  const ec = new EC('secp256k1')
  const keypair = ec.genKeyPair()

  beforeEach(() => {
    // catch everything on stdout>
    // and put it in result
    result = []

    // this hides console.log calls
    jest.spyOn(process.stdout, 'write').mockImplementation(val => {
      result.push(val)
    })
  })

  afterEach(() => jest.restoreAllMocks())

  it('Login', async () => {
    // Start logged out
    await logout()

    const pubPoint = keypair.getPublic()
    const fakePublicKey = pubPoint.encode('hex')
    const fakePrivateKey = keypair.getPrivate('hex')

    // overwrite static function getKeypair in order to know the keypair
    LoginCommand.getKeypair = () => {
      return {
        publicKey: fakePublicKey,
        privateKey: fakePrivateKey,
      }
    }

    const simulateUserConfirm = async () => {
      // simulate the user login to generate a file with the session
      // encrypted with the key above.
      await login()
      const userSession = await getSessionWithConfig()

      await putFile({
        path: `${fakePublicKey}.json`,
        content: JSON.stringify(userSession.store.getSessionData()),
        encrypt: fakePublicKey,
      })

      await logout()
    }

    await simulateUserConfirm()
    await runCommand(LoginCommand, ['samuelroy.id.blockstack'])

    expect(open).toHaveBeenCalledTimes(1)

    const logged = result.find(log => {
      return log.indexOf(`You can logout with`) > -1
    })
    expect(logged).toBeDefined()
  }, 10000)
})
