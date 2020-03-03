const nock = require('nock')
const path = require('path')

const fs = require('fs')
const rimraf = require('rimraf')
const { write, read } = require('../../lib/cliStorage')
const { privateKey } = require('./keypair')
const InitCommand = require('../../commands/init')
const { login, logout, runCommand } = require('./helpers')

const fsp = fs.promises

const mockGaiaToLocalFileSystem = () => {
  nock('https://hub.blockstack.org')
    .persist()
    .post(/store\/.*/)
    .reply((uri, body) => {
      uri = uri.replace(/\/store\/12ENwmKf2wn5AS63i8cSTyhBvTXK4EXB1y/, '')
      write({
        path: path.join(__dirname, '../hub'),
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
        path: path.join(__dirname, '../hub'),
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
  await fsp.mkdir(path.join(__dirname, '../hub'))
  await fsp.mkdir(path.join(__dirname, '../local/'))
  await fsp.writeFile(path.join(__dirname, '../local/foo.txt'), 'foo bar')
  await login()
  await runCommand(InitCommand, ['unit-test-project'])
}

module.exports = {
  mockGaiaToLocalFileSystem,
  createDescriptor,
  prepareEnvironment,
}
