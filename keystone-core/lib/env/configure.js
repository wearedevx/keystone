const debug = require('debug')('keystone:core:env')

const { getPath } = require('../descriptor-path')
const {
  updateDescriptorForMembers,
  getMembers,
  removeDescriptorForMembers,
  getLatestEnvDescriptor,
  getLatestProjectDescriptor,
} = require('../descriptor')

const { createMembersDescriptor } = require('../member')

const createEnv = async (userSession, { env, projectDescriptor }) => {
  const envContent = {
    name: env,
    files: [],
  }

  try {
    const envMembersDescriptor = await createMembersDescriptor(userSession, {
      project: projectDescriptor.name,
      env,
    })
    // create env descriptor and upload for every project members
    const envDescriptor = await updateDescriptorForMembers(userSession, {
      project: projectDescriptor.name,
      name: env,
      type: 'env',
      env,
      membersDescriptor: envMembersDescriptor,
      content: envContent,
    })

    return envDescriptor
  } catch (err) {
    console.error(err)
    throw new Error('Project file could not have been connected')
  }
}

const getFilesFromEnv = async (userSession, { project, env }) => {
  const { username } = userSession.loadUserData()

  const envDescriptor = await getLatestEnvDescriptor(userSession, {
    project,
    type: 'env',
    env,
    blockstackId: username,
  })

  return envDescriptor.content.files
}

const removeFilesFromEnv = async (
  userSession,
  { projectDescriptor, env, absoluteProjectPath, members }
) => {
  debug('removeFilesFromEnv')

  const files = await getFilesFromEnv(userSession, {
    project: projectDescriptor.name,
    env,
  })

  if (!members) {
    members = await getMembers(userSession, {
      project: projectDescriptor.name,
      env,
    })
  }

  return Promise.all(
    files.map(file =>
      removeDescriptorForMembers(userSession, {
        project: projectDescriptor.name,
        env,
        descriptorPath: file.name,
        absoluteProjectPath,
        members,
        type: 'file',
      })
    )
  )
}

const removeEnvFiles = async (
  userSession,
  { project, env, absoluteProjectPath }
) => {
  debug('removeEnvFiles', env)

  const { username } = userSession.loadUserData()

  // Check if env exists.
  const projectDescriptor = await getLatestProjectDescriptor(userSession, {
    project,
    type: 'project',
    blockstackId: username,
  })

  // If not, throw an error.
  if (!projectDescriptor.content.env.includes(env)) {
    throw new Error(`Env ${env} does not exist.`)
  }

  const members = await getMembers(userSession, { project, env })

  try {
    await removeFilesFromEnv(userSession, {
      projectDescriptor,
      env,
      absoluteProjectPath,
      members,
    })
  } catch (err) {
    console.error(err)
  }

  // // // Remove env descriptor
  const envPath = getPath({ project, env, type: 'env' })

  try {
    await removeDescriptorForMembers(userSession, {
      project,
      env,
      members,
      type: 'env',
      descriptorPath: envPath,
    })
  } catch (err) {
    console.error(err)
  }
}

const setMembersToEnvs = ({ envsMembers, members, role, env }) => {
  const selectedEnv = envsMembers[env]
  Object.keys(selectedEnv).forEach(group => {
    selectedEnv[group] = selectedEnv[group].filter(
      member =>
        !members.find(
          newMember => newMember.blockstack_id === member.blockstack_id
        )
    )
  })
  // Prevent duplicates
  selectedEnv[role] = members.reduce((dedupMembers, member) => {
    const found = dedupMembers.find(
      m => m.blockstack_id === member.blockstack_id
    )
    if (found) {
      return dedupMembers
    }
    dedupMembers.push(member)
    return dedupMembers
  }, [])

  return { ...envsMembers, [env]: selectedEnv }
}

const isOneOrMoreAdmin = envsMembers => {
  return (
    Object.keys(envsMembers).filter(env => envsMembers[env].admins.length === 0)
      .length === 0
  )
}

module.exports = {
  createEnv,
  removeEnvFiles,
  isOneOrMoreAdmin,
  setMembersToEnvs,
}
