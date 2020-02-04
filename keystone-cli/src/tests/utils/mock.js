const debug = require('debug')('keystone:core:file')
const fs = require('fs')
const nock = require('nock')
const pth = require('path')
const { SHARED_MEMBER, PUBKEY } = require('@keystone.sh/core/lib/constants')
const fsp = fs.promises

const gaiaFile = require('@keystone.sh/core/lib/file/gaia')

gaiaFile.readFileFromGaia = async (
  userSession,
  { path, json = true, origin, decrypt }
) => {
  const { username } = userSession.loadUserData()
  if (decrypt && userSession.sharedPrivateKey) {
    decrypt = userSession.sharedPrivateKey
  }
  console.log(path)
  path = path.replace(/\//g, '|')
  console.log('READ FILE ', path)

  path = `${username}--${path}`
  path = pth.join(__dirname, '..', pth.sep, 'hub', path)

  try {
    const fetchedFile = await fsp.readFile(path, 'utf-8')
    if (json) {
      const fetchedFileJSON = JSON.parse(fetchedFile)
      return fetchedFileJSON
    }

    return fetchedFile
  } catch (error) {
    // we can't retrieve the file from remote
    // make it like a 404 not found
    console.error(error)
    return null
  }
}

gaiaFile.writeFileToGaia = async (userSession, { path, content, encrypt }) => {
  const { username } = userSession.loadUserData()
  path = path.replace(/\//g, '|')
  console.log('WRITE FILE ', path)
  path = `${username}--${path}`
  path = pth.join(__dirname, '..', pth.sep, 'hub', path)
  await fsp.writeFile(path, content)
  return content
}

gaiaFile.getPubkey = async (
  userSession,
  { blockstack_id: blockstackId, publicKey }
) => {
  if (publicKey) return publicKey

  if (new RegExp(SHARED_MEMBER).test(blockstackId))
    return blockstackId.split('-')[1]
  const pubkeyFile = await gaiaFile.readFileFromGaia(userSession, {
    decrypt: false,
    path: `${PUBKEY}`,
    origin: blockstackId,
    json: false,
    verify: true,
  })
  if (pubkeyFile) {
    return pubkeyFile
  }
  throw new Error(
    `Keystone public application key not found on ${blockstackId}`
  )
}

gaiaFile.listFilesFromGaia = async userSession => {
  return fs.readdirSync(pth.join(__dirname, '..', pth.sep, 'hub'))
}
module.exports = gaiaFile
