const fs = require('fs')
const NewCommand = require('../commands/new')
const RemoveCommand = require('../commands/remove')
const { login, logout, runCommand } = require('./helpers')

describe('New command', () => {
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

  it('new - Not connected', async () => {
    await logout()
    await runCommand(NewCommand, ['test-project'])

    const createdProject = result.find(log => {
      return /You're not connected, please sign in first/g.test(log)
    })
    expect(createdProject).toBeDefined()
  })

  it('should create a new project', async () => {
    await login()
    await runCommand(NewCommand, ['new-project'])

    const createdProject = result.find(log => {
      return /.*Project .*new-project.* successfully created.*/g.test(log)
    })
    expect(createdProject).toBeDefined()

    await runCommand(RemoveCommand, [
      `--project=${createdProject.replace(/.[[0-9]+m/g, '').split(' ')[2]}`,
    ])
  }, 20000)
  it('should not  create a project because no name provided', async () => {
    await login()
    let noArgumentErr
    try {
      await runCommand(NewCommand, [])
    } catch (err) {
      noArgumentErr = err
    }

    expect(noArgumentErr).toBeDefined()
  })
})
