const EC = require('elliptic').ec
const { addMember } = require('../../member')
const { ROLES, SHARED_MEMBER } = require('../../constants')
const pull = require('../pull')

const ec = new EC('secp256k1')

const newShare = async (userSession, { project, env }) => {
  const { username } = userSession.loadUserData()
  const keypair = ec.genKeyPair()
  const pubPoint = keypair.getPublic()

  const userKeypair = {
    publicKey: pubPoint.encode('hex'),
    privateKey: keypair.getPrivate('hex'),
  }
  await addMember(userSession, {
    project,
    env,
    member: SHARED_MEMBER,
    role: ROLES.READERS,
    publicKey: userKeypair.publicKey,
  })
  await pull(userSession, {
    project,
    env,
    origin: username,
    absoluteProjectPath: './',
    force: true,
  })
  return userKeypair
}

module.exports = { newShare }
