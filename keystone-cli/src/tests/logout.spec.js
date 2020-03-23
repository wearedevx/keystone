process.env.SESSION_FILENAME = 'session-test1.json'

require('./utils/mock')
const fs = require('fs')

const { prepareEnvironment } = require('./utils')

jest.mock('../lib/blockstackLoader')
jest.mock('../lib/commands')

const { runCommand } = require('./utils/helpers')
const LogoutCommand = require('../commands/logout')

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
    await prepareEnvironment()

    await runCommand(LogoutCommand)

    const logged = result.find(log => {
      return log.indexOf(`Sign out from`) > -1
    })
    expect(logged).toBeDefined()
  })
})
