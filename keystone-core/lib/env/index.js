const debug = require('debug')('keystone:core:env')
const fs = require('fs')
const path = require('path')

const { getPath } = require('../descriptor-path')
const {
  updateDescriptorForMembers,
  getMembers,
  removeDescriptorForMembers,
  getLatestProjectDescriptor,
  getLatestEnvDescriptor,
} = require('../descriptor')
const { KEYSTONE_CONFIG_PATH } = require('../constants')

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
  const projectDescriptor = await getLatestEnvDescriptor(userSession, {
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
  Object.keys(selectedEnv).forEach(r => {
    selectedEnv[r] = selectedEnv[r].filter(
      member => !members.find(x => x.blockstack_id === member.blockstack_id)
    )
  })
  selectedEnv[role] = members

  return { ...envsMembers, [env]: selectedEnv }
}

const isOneOrMoreAdmin = envsMembers => {
  return (
    Object.keys(envsMembers).filter(env => envsMembers[env].admins.length === 0)
      .length === 0
  )
}

const checkoutEnv = async (
  userSession,
  { project, env, absoluteProjectPath }
) => {
  const projectDescriptor = await getLatestProjectDescriptor(userSession, {
    project,
    type: 'project',
  })

  // Retrieve updated project descriptor
  console.log('projectDescriptor', projectDescriptor.content.members)
  console.log('projectDescriptor', projectDescriptor.content.env.includes(env))

  const envFound = projectDescriptor.content.env.find(
    envObject => envObject === env
  )

  if (envFound) {
    await getLatestEnvDescriptor(userSession, {
      project,
      env,
      type: 'env',
    })
    const configFile = JSON.parse(
      fs.readFileSync(path.join(absoluteProjectPath, KEYSTONE_CONFIG_PATH))
    )
    configFile.env = env
    fs.writeFileSync('.ksconfig', JSON.stringify(configFile))
    return configFile
  }
  throw new Error(`The environment ${env} is not defined in this project`)
}

module.exports = {
  createEnv,

  removeEnvFiles,
  isOneOrMoreAdmin,
  setMembersToEnvs,
  checkoutEnv,
}
