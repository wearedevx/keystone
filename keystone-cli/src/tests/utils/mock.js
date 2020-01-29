const debug = require('debug')('keystone:core:file')
const fs = require('fs')
const nock = require('nock')
const pth = require('path')
const fsp = fs.promises

const gaiaFile = require('@keystone.sh/core/lib/file')

gaiaFile.readFileFromGaia = async (
  userSession,
  { path, json = true, origin, decrypt }
) => {
  console.log("ENTER THE MOCKED 'GET FILE'")
  const { username } = userSession.loadUserData()
  if (decrypt && userSession.sharedPrivateKey) {
    decrypt = userSession.sharedPrivateKey
  }

  debug('')

  debug(`Get file from gaia ${path} from ${origin || 'self'}`)

  path = `${username}--${path}`
  path = pth.join(__dirname, '../hub', path)

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
  debug('')
  debug(`Write file on gaia to ${path}`)
  path = `${username}--${path}`
  path = pth.join(__dirname, '../hub', path)
  await fsp.writeFile(path, content)
  return content
}

module.exports = gaiaFile
