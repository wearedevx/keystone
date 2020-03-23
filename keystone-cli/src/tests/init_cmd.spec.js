require('./utils/mock')

jest.mock('../lib/blockstackLoader')
jest.mock('../lib/commands')

const { stdin } = require('mock-stdin')
const fs = require('fs')

const { prepareEnvironment } = require('./utils')
const InitCommand = require('../commands/init')
const ProjectListCommand = require('../commands/list')
const ProjectRmCommand = require('../commands/project/rm')
const { runCommand } = require('./utils/helpers')

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
    await runCommand(ProjectListCommand)
    let existingProject = result.find(
      log => log.indexOf(`> ${PROJECT_NAME}`) > -1
    )

    // If project already exist, remove it from hub
    if (existingProject) {
      existingProject = existingProject.replace(/.[[0-9]+m/g, '')
      existingProject = existingProject.match(
        /unit-test-project\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/ // match project name followed by uuid
      )
      const sendKeystrokes = async () => {
        io.send(keys.enter)
        io.send(existingProject)
      }
      setTimeout(() => sendKeystrokes().then(), 1000)
      await runCommand(ProjectRmCommand, [existingProject])
    }

    const sendKeystrokes = async () => {
      io.send(keys.enter)
    }
    setTimeout(() => sendKeystrokes().then(), 500)
    await runCommand(InitCommand, [PROJECT_NAME])
    const createdProject = result.find(log =>
      /.* successfully created/g.test(log)
    )
    expect(createdProject).toBeDefined()
  }, 20000)
})
