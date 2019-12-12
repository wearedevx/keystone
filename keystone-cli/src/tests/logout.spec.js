const LogoutCommand = require('../commands/logout')
const { login, logout, runCommand } = require('./helpers')

describe('Logout Command', () => {
  let result

  beforeEach(() => {
    // catch everything on stdout>
    // and put it in result
    result = []

    // /!\ this hides console.log calls
    jest.spyOn(process.stdout, 'write').mockImplementation(val => {
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
