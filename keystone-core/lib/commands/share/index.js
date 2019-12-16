const EC = require('elliptic').ec
const { addMember } = require('../../member')
const { ROLES } = require('../../constants')

const ec = new EC('secp256k1')

const newShare = (userSession, { project }) => {
  const keypair = ec.genKeyPair()
  const pubPoint = keypair.getPublic()

  const userKeypair = {
    publicKey: pubPoint.encode('hex'),
    privateKey: keypair.getPrivate('hex'),
  }
  return addMember(userSession, {
    project,
    member: `shared-${userKeypair.publicKey}`,
    role: ROLES.SHARES,
  })
}

module.exports = { newShare }
