const ListCommand = require('../commands/list')
const { login, logout, runCommand } = require('./helpers')

describe('List command', () => {
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

  it('should not list because user is not connected', async () => {
    await logout()
    await runCommand(ListCommand, [])
    expect(result[result.length - 1]).toContain(
      `You're not connected, please sign in first`
    )
  })

  it('should list all the projects', async () => {
    await login()
    await runCommand(ListCommand, [])
    const projectFound = result.find(log => {
      return /Projects: .* found/g.test(log)
    })
    expect(projectFound).toBeDefined()
  })

  it('should list all files in specific project', async () => {
    await login()

    await runCommand(ListCommand, ['--project=keystone'])

    const filesFound = result.find(log => {
      return /Files: .*/g.test(log)
    })
    expect(filesFound).toBeDefined()
  })
})
