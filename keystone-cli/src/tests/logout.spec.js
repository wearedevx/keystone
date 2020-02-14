process.env.SESSION_FILENAME = 'session-test1.json'

require('./utils/mock')
const { prepareEnvironment } = require('./utils')

const fs = require('fs')

jest.mock('../lib/blockstackLoader')
jest.mock('../lib/commands')

const { login, logout, runCommand } = require('./utils/helpers')
const LogoutCommand = require('../commands/logout')
const { SESSION_FILENAME } = require('@keystone.sh/core/lib/constants')
describe('Logout Command', () => {
  let result

  beforeEach(() => {
    // catch everything on stdout
    // and put it in result
    result = []

    // /!\ this hides console.log calls
    jest.spyOn(process.stdout, 'write').mockImplementation(val => {
      fs.appendFile('unit-test.log', val)
      result.push(val)
    })
  })

  afterEach(() => jest.restoreAllMocks())

  it('Logout', async () => {
    // Start with a session signed in.
    await login()
    await runCommand(LogoutCommand)

    const logged = result.find(log => {
      return log.indexOf(`Sign out from`) > -1
    })
    expect(logged).toBeDefined()
  })

  it(`Can't logout if not signed in`, async () => {
    // We start logged out
    await logout()

    await runCommand(LogoutCommand)

    const logged = result.find(log => {
      return log.indexOf(`You're not connected, please sign in first`) > -1
    })
    expect(logged).toBeDefined()
  })
})
