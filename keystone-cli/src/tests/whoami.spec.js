const WhoAmICommand = require('../commands/whoami')
const { login, logout, runCommand } = require('./helpers')

describe('Who Am I Command', () => {
  let result

  beforeEach(() => {
    // catch everything on stdout
    // and put it in result
    result = []
    jest
      .spyOn(process.stdout, 'write')
      .mockImplementation(val => result.push(val))
  })

  afterEach(() => jest.restoreAllMocks())

  it('Who am I - Not connected', async () => {
    await logout()
    await runCommand(WhoAmICommand, [])
    expect(result[result.length - 1]).toContain(
      `You're not connected, please sign in first`
    )
  })

  it('Who am I - connected', async () => {
    await login()
    await runCommand(WhoAmICommand, [])
    expect(result[result.length - 1]).toContain(`You can logout with`)
  })
})
