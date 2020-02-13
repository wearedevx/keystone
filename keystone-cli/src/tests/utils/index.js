const nock = require('nock')
const path = require('path')

const fs = require('fs')
const rimraf = require('rimraf')
const { write, read } = require('../../lib/cliStorage')
const { privateKey } = require('./keypair')
const InitCommand = require('../../commands/init')
const { login, logout, runCommand } = require('./helpers')

const mockGaiaToLocalFileSystem = () => {
  nock('https://hub.blockstack.org')
    .persist()
    .post(/store\/.*/)
    .reply((uri, body) => {
      uri = uri.replace(/\/store\/12ENwmKf2wn5AS63i8cSTyhBvTXK4EXB1y/, '')
      write({
        path: '/home/abigael/code/keystone/keystone-cli/src/tests/hub',
        filename: uri,
        content: JSON.stringify(body),
      })

      return [200, {}]
    })

  nock('https://hub.blockstack.org')
    .persist()
    .get(/store\/.*/)
    .reply(async (uri, body) => {
      uri = uri.replace(/\/store\/12ENwmKf2wn5AS63i8cSTyhBvTXK4EXB1y/, '')
      const { decryptECIES } = require('blockstack/lib/encryption/ec')

      const uncryptedData = await read({
        path: '/home/abigael/code/keystone/keystone-cli/src/tests/hub',
        filename: uri,
      })

      const data = await decryptECIES(privateKey, uncryptedData)

      return [200, data]
    })
}

const createDescriptor = ({
  name = 'foo.txt',
  content = 'foo bar',
  type = 'file',
  path = '',
  checksum = '',
  version = 1,
  history = [],
}) => {
  return { name, content, path, checksum, type, version, history }
}

const prepareEnvironment = async () => {
  rimraf.sync(path.join(__dirname, '../hub/'))
  rimraf.sync(path.join(__dirname, '../local/'))
  await login()
  await runCommand(InitCommand, ['unit-test-project'])
}

module.exports = {
  mockGaiaToLocalFileSystem,
  createDescriptor,
  prepareEnvironment,
}
