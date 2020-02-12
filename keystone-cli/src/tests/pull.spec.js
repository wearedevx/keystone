require('./utils/mock')
const fs = require('fs')
const uuid = require('uuid/v4')
const pathUtil = require('path')
const { writeFileToDisk } = require('@keystone.sh/core/lib/file')
const PullCommand = require('../commands/pull')
const PushCommand = require('../commands/push')
const DeleteCommand = require('../commands/delete')

const { login, logout, runCommand, addMemberToEnv } = require('./utils/helpers')
const { createDescriptor } = require('./utils')
jest.mock('../lib/blockstackLoader')
jest.mock('../lib/commands')

const fsp = fs.promises

describe('Push Command', () => {
  let result

  beforeEach(() => {
    // catch everything on stdout
    // and put it in result
    result = []
    jest.spyOn(process.stdout, 'write').mockImplementation(val => {
      fs.appendFile('unit-test.log', val)
      result.push(val)
    })
  })

  afterEach(() => jest.restoreAllMocks())

  it('should not update any files because same version', async () => {
    await login()

    await runCommand(PullCommand)

    const pulledFile = result.find(
      log => log.indexOf('You are already up to date') > -1
    )
    expect(pulledFile).toBeDefined()
  })

  fit('should pull a file pushed by another user', async () => {
    await login()
    //add other member ot the environment
    await addMemberToEnv({ username: 'keystone_test2.id.blockstack' })

    // Connect to the user two
    await login(2)
    await runCommand(PullCommand)

    const uid = uuid()
    const fileDescriptor = createDescriptor({ content: uid })
    await writeFileToDisk(fileDescriptor)
    await runCommand(PushCommand, [pathUtil.join(fileDescriptor.name)])
    await logout()

    await login()
    await runCommand(DeleteCommand, ['foo.txt'])

    await runCommand(PullCommand)

    const content = (
      await fsp.readFile(pathUtil.join(__dirname, 'local/foo.txt'))
    ).toString()

    expect(uid).toEqual(content)
  })

  // fit('should update a file from because newer version on storage', async () => {
  //   await login()
  //   let fileToChange
  //   let envDescriptorToChange
  //   const files = fs.readdirSync(path.join(__dirname, './hub'))

  //   files.forEach(file => {
  //     if (file.indexOf('foo.txt') > -1) {
  //       fileToChange = path.join(__dirname, './hub/', file)
  //     }
  //     if (file.indexOf('default|') > -1) {
  //       envDescriptorToChange = path.join(__dirname, './hub/', file)
  //     }
  //   })

  //   const fileDescriptor = JSON.parse(fs.readFileSync(fileToChange))
  //   const envDescriptor = JSON.parse(fs.readFileSync())
  //   const newFileDescriptor = {
  //     ...fileDescriptor,
  //     version: fileDescriptor.version + 1,
  //     content: `${fileDescriptor.content} quu`,
  //   }

  //   console.log('FILE DESCRIPTOR', newFileDescriptor)
  //   return

  //   fs.writeFileSync(fileToChange, JSON.stringify(newFileDescriptor))

  //   await runCommand(PullCommand)

  //   const pulledFile = result.find(
  //     log => log.indexOf('You are already up to date') > -1
  //   )

  //   expect(pulledFile).toBeDefined()
  // })
})
