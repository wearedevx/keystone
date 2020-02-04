const nock = require('nock')
const path = require('path')
const { write, read } = require('../../lib/cliStorage')
const { privateKey } = require('./keypair')

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
      console.log('je suis charlie')
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

module.exports = {
  mockGaiaToLocalFileSystem,
}
