require('./utils/mock')
const { prepareEnvironment } = require('./utils')

jest.mock('../lib/blockstackLoader')
jest.mock('../lib/commands')

const { stdin } = require('mock-stdin')
const fs = require('fs')
const InitCommand = require('../commands/init')
const ListCommand = require('../commands/list')
const DeleteCommand = require('../commands/delete')
const { login, logout, runCommand } = require('./utils/helpers')

// Key codes
const keys = {
  up: '\x1B\x5B\x41',
  down: '\x1B\x5B\x42',
  enter: '\x0D',
  space: '\x20',
}

const PROJECT_NAME = 'unit-test-project'

describe('Init Command', () => {
  let result
  let io

  beforeEach(() => {
    // catch everything on stdout>
    // and put it in result
    result = []

    // /!\ this hides console.log calls
    jest.spyOn(process.stdout, 'write').mockImplementation(val => {
      fs.appendFile('unit-test.log', val)
      result.push(val)
    })

    // helper for sending keystrokes
    io = stdin()
  })

  afterEach(() => {
    jest.restoreAllMocks()
    io.restore()
  })

  it('should create a new config and a new project', async () => {
    await prepareEnvironment()

    // remove any existing config
    if (fs.existsSync('.ksconfig')) {
      fs.unlinkSync('.ksconfig')
    }
    // remove the project if already exists
    await runCommand(ListCommand, ['projects'])
    let existingProject = result.find(log => log.indexOf(PROJECT_NAME) > -1)
    // match project name followed by uuid
    if (existingProject) {
      existingProject = existingProject
        .replace(/.[[0-9]+m/g, '')
        .match(
          /unit-test-project\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/
        )
      const sendKeystrokes = async () => {
        io.send(keys.enter)
      }
      setTimeout(() => sendKeystrokes().then(), 500)
      await runCommand(DeleteCommand, [`--project=${existingProject}`])
    }
    await runCommand(InitCommand, [PROJECT_NAME])
    const createdProject = result.find(log =>
      /.* successfully created/g.test(log)
    )
    expect(createdProject).toBeDefined()
  }, 20000)

  // it('should overwrite an existing config if user confirms', async () => {
  //   // Start with a session signed in.
  //   await login()
  //
  //   const existingConfig = fs.readFileSync('.ksconfig')
  //
  //   // Send Keystroke to confirm overwriting the config file
  //   const sendKeystrokes = async () => {
  //     io.send(keys.enter)
  //   }
  //   setTimeout(() => sendKeystrokes().then(), 500)
  //
  //   await runCommand(InitCommand, [PROJECT_NAME])
  //
  //   const configFileCreated = result.find(log => {
  //     return log.indexOf(`.ksconfig file created`) > -1
  //   })
  //
  //   expect(configFileCreated).toBeDefined()
  //
  //   // Delete the project
  //   // Set config file as before
  //   const createdProject = result.find(log =>
  //     /.*Project .* successfully created/g.test(log)
  //   )
  //
  //   setTimeout(() => sendKeystrokes().then(), 500)
  //   await runCommand(RemoveCommand, [
  //     `--project=${createdProject.replace(/.[[0-9]+m/g, '').split(' ')[2]}`,
  //   ])
  //
  //   fs.writeFile('.ksconfig', existingConfig)
  // }, 20000)
  //
  // it(`should not initialize a project if the user is not logged in`, async () => {
  //   // We start logged out
  //   await logout()
  //
  //   await runCommand(InitCommand, [PROJECT_NAME])
  //
  //   const needToBeLogged = result.find(log => {
  //     return log.indexOf(`You're not connected, please sign in first`) > -1
  //   })
  //   expect(needToBeLogged).toBeDefined()
  // })
})
