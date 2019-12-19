const EC = require('elliptic').ec
const { addMember } = require('../../member')
const { ROLES, SHARED_MEMBER } = require('../../constants')
const { getPath } = require('../../descriptor-path')
const { getLatestDescriptorByPath } = require('../../descriptor')
const { writeFileToDisk } = require('../../file/disk')

const { extractMembersByRole } = require('../../descriptor')

const ec = new EC('secp256k1')

const getLastVersionOfMemberDescriptor = async (
  userSession,
  { project, env, members, stable = false }
) => {
  const memberPath = getPath({
    env,
    blockstackId: SHARED_MEMBER,
    project,
    type: 'members',
  })

  const memberDescriptor = await getLatestDescriptorByPath(
    userSession,
    {
      descriptorPath: memberPath,
      members,
    },
    stable
  )

  // If the is admin in fetched descriptor that we didn't know before
  if (
    memberDescriptor[0].content[ROLES.ADMINS].find(
      m => !members.find(me => me.blockstack_id === m.blockstack_id)
    )
  ) {
    return getLastVersionOfMemberDescriptor(userSession, {
      project,
      env,
      members: memberDescriptor[0].content[ROLES.ADMINS],
      stable,
    })
  }
  return memberDescriptor[0]
}

const pullShared = async (
  userSession,
  { project, env, origins: members, absoluteProjectPath }
) => {
  const memberDescriptor = await getLastVersionOfMemberDescriptor(userSession, {
    project,
    env,
    members,
    stable: true,
  })
  const envPath = getPath({
    env,
    blockstackId: SHARED_MEMBER,
    project,
    type: 'env',
  })

  const membersToRetrieveFiles = extractMembersByRole(memberDescriptor, [
    ROLES.ADMINS,
    ROLES.CONTRIBUTORS,
  ])

  const envDescriptor = await getLatestDescriptorByPath(
    userSession,
    {
      descriptorPath: envPath,
      members: membersToRetrieveFiles,
    },
    true
  )

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
        members: membersToRetrieveFiles,
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

  const memberDescriptor = await addMember(userSession, {
    project,
    env,
    member: SHARED_MEMBER,
    role: ROLES.READERS,
    publicKey: userKeypair.publicKey,
  })

  return { privateKey: userKeypair.privateKey, memberDescriptor }
}

module.exports = { newShare, pullShared }
