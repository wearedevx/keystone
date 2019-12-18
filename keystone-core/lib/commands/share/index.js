const EC = require('elliptic').ec
const { addMember } = require('../../member')
const { ROLES, SHARED_MEMBER } = require('../../constants')
const { getPath } = require('../../descriptor-path')
const { getLatestDescriptorByPath } = require('../../descriptor')
const { writeFileToDisk } = require('../../file/disk')
const ec = new EC('secp256k1')

const pullShared = async (
  userSession,
  { project, env, origin, absoluteProjectPath }
) => {
  const envPath = getPath({
    env,
    blockstackId: SHARED_MEMBER,
    project,
    type: 'env',
  })

  const envDescriptor = await getLatestDescriptorByPath(
    userSession,
    { descriptorPath: envPath, members: [{ blockstack_id: origin }] },
    true
  )
  console.log('TCL: pullShared -> envDescriptor', envDescriptor.content)

  const files = await Promise.all(
    envDescriptor[0].content.files.map(async ({ name: filename }) => {
      const path = getPath({
        env,
        blockstackId: SHARED_MEMBER,
        project,
        filename,
        type: 'file',
      })

      const fileDescriptor = await getLatestDescriptorByPath(userSession, {
        descriptorPath: path,
        members: [{ blockstack_id: origin }],
      })

      writeFileToDisk(fileDescriptor, absoluteProjectPath)
    })
  )
}

const newShare = async (userSession, { project, env }) => {
  const keypair = ec.genKeyPair()
  const pubPoint = keypair.getPublic()

  const userKeypair = {
    publicKey: pubPoint.encode('hex'),
    privateKey: keypair.getPrivate('hex'),
  }

  console.log(userKeypair)

  const memberDescriptor = await addMember(userSession, {
    project,
    env,
    member: SHARED_MEMBER,
    role: ROLES.READERS,
    publicKey: userKeypair.publicKey,
  })

  return userKeypair
}

module.exports = { newShare, pullShared }
