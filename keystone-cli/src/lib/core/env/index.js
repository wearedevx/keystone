const fs = require('fs')
const hash = require('object-hash')

const { deepCopy } = require('../../utils')

const { uploadFile, getPath, deleteFiles } = require('../file')

const {
  uploadFileInProject,
  getEnvsDescriptor,
  getProjectDescriptor,
  getAllMembers,
} = require('../project')

const isSameFileInEnv = ({ envDescriptor, filename }) => {
  return envDescriptor.content.files.find(f => f.name === filename)
}

const getAdminAndContributorIds = envDescriptor => {
  return [
    ...envDescriptor.content.members.contributors.map(x => x.blockstack_id),
    ...envDescriptor.content.members.admins.map(x => x.blockstack_id),
  ]
}
const getEnvMembers = envDescriptor => {
  return [
    ...envDescriptor.content.members.contributors.map(x => x.blockstack_id),
    ...envDescriptor.content.members.admins.map(x => x.blockstack_id),
    ...envDescriptor.content.members.readers.map(x => x.blockstack_id),
  ]
}

const uploadFileToEnv = async (
  userSession,
  { projectDescriptor, fileDescriptor, env, envDescriptor }
) => {
  await uploadFileInProject(userSession, {
    projectDescriptor,
    fileDescriptor,
    type: 'file',
    env,
    envDescriptor,
    versioning: true,
    members: getEnvMembers(envDescriptor),
  })

  // const filename = extractFileFromPath(fileDescriptor.name)
  const newFile = {
    // localPath: fileDescriptor.localPath,
    name: fileDescriptor.name,
  }

  if (isSameFileInEnv({ envDescriptor, filename: fileDescriptor.name }))
    return envDescriptor

  // if the file is already present but not the same (new tags or tags removed)
  // filter out
  // envDescriptor.content.files = envDescriptor.content.files.filter(
  //   file => file.name !== filename
  // )

  envDescriptor.content.files.push(newFile)

  await uploadFileInProject(userSession, {
    projectDescriptor,
    fileDescriptor: envDescriptor,
    type: 'env',
    env,
    envDescriptor,
    versioning: true,
    members: getEnvMembers(envDescriptor),
  })

  return {
    ...envDescriptor,
    content: {
      ...envDescriptor.content,
      files: [...envDescriptor.content.files, newFile],
    },
  }
}

const updateProjectEnv = async (
  userSession,
  { projectDescriptor, env, path }
) => {
  // try {
  // const projectDescriptor = await getProjectDescriptor(userSession, {
  //   project,
  //   type: 'project',
  // })

  const projectChecksumLast = hash(projectDescriptor)

  if (projectDescriptor.content.envs.find(x => x.name === env)) {
    throw new Error(`The environement "${env}" is already set for this project`)
  }

  projectDescriptor.content.envs.push({ name: env, path })

  const projectChecksumNew = hash(projectDescriptor)

  // try {
  if (projectChecksumLast !== projectChecksumNew) {
    // should update

    await uploadFileInProject(userSession, {
      projectDescriptor,
      fileDescriptor: projectDescriptor,

      // versioning: true,
      type: 'project',
      members: getAllMembers(projectDescriptor),
      env,
      versioning: true,
    })
    // await userSession.putFile(
    //   projectDescriptor.path,
    //   JSON.stringify(projectDescriptor),
    //   { encrypt: true }
    // )
  }
  //   } catch (error) {
  //     throw error
  //   }
  // } catch (err) {
  //   throw err
  // }
}

const createEnv = async (userSession, { env, projectDescriptor }) => {
  const { username, email } = userSession.loadUserData()

  const path = getPath({
    project: projectDescriptor.content.name,
    env,
    blockstack_id: username,
    type: 'env',
  })
  const envContent = {
    name: env,
    // {path, name}
    files: [],
    members: {
      admins: [
        {
          blockstack_id: username,
          email,
        },
      ],
      contributors: [],
      readers: [],
    },
  }
  try {
    await updateProjectEnv(userSession, {
      projectDescriptor,
      path,
      env,
      members: envContent.members,
    })
  } catch (err) {
    console.error(err)
    throw new Error('Project file could not have been connected')
  }

  const envDescriptor = {
    path,
    checksum: hash(JSON.stringify(envContent)),
    content: envContent,
    version: 1,
    author: blockstack_id,
    history: {},
  }

  try {
    await uploadFileInProject(userSession, {
      path,
      content: JSON.stringify(envDescriptor.content),
      envDescriptor,
      projectDescriptor,
      fileDescriptor: envDescriptor,
      file: envDescriptor.content.name,
      versioning: true,
      type: 'env',
      env,
    })

    return envDescriptor
  } catch (error) {
    console.log("Couldn't save the env", error)
  }
}

const removeEnv = async (userSession, { name, envsDescriptor }) => {
  await deleteFiles(userSession, { env: [name], envsDescriptor })
  await removeEnvFile(userSession, { name })
  await removeEnvFromProject(userSession, { name })
}

const checkoutEnv = async (userSession, { env, projectDescriptor }) => {
  // Retrieve updated project descriptor
  console.log('projectDescriptor', projectDescriptor.content.members)
  console.log('projectDescriptor', projectDescriptor.content.envs.includes(env))

  const envFound = projectDescriptor.content.envs.find(
    envObject => envObject.name === env
  )

  if (envFound) {
    const configFile = JSON.parse(fs.readFileSync('.ksconfig'))
    configFile.env = env
    fs.writeFileSync('.ksconfig', JSON.stringify(configFile))
    return configFile
  }
  throw new Error(`The environment ${env} is not defined in this project`)
}

const getEnv = () => {
  try {
    const { env } = JSON.parse(fs.readFileSync('.ksconfig'))
    return env
  } catch (err) {
    throw err
  }
}

const updateMembersRoles = async (
  userSession,
  { projectDescriptor, envsDescriptor }
) => {
  try {
    await Promise.all(
      envsDescriptor.map(async descriptor => {
        try {
          // const envChecksum = hash(descriptor)

          // TODO update env file for all members
          // await copyFilesForMembers(userSession, {
          //   projectDescriptor,
          // })

          // project descriptor has been updated
          await uploadFile(userSession, {
            project: projectDescriptor.content.name,
            members: getAdminAndContributorIds(descriptor),
            fileDescriptor: descriptor,
            versioning: true,
            type: 'env',
            env: descriptor.content.name,
            path: descriptor.path,
          })
        } catch (err) {
          throw err
        }
      })
    )
  } catch (error) {
    console.log(error)
    throw error
  }
}

const removeEnvFromProject = async (userSession, { name }) => {
  const projectDescriptor = await getProjectDescriptor(userSession, {})
  projectDescriptor.content.envs = projectDescriptor.content.envs.filter(
    env => env.name !== name
  )
  try {
    uploadFile(userSession, {
      path: projectDescriptor.path,
      fileDescriptor: projectDescriptor,
      type: 'project',
      file: projectDescriptor.content.name,
    })
  } catch (err) {
    throw err
  }
}

const deleteFilesFromEnv = async (
  userSession,
  { projectDescriptor, files, envDescriptor, env }
) => {
  const envDescriptorCloned = deepCopy(envDescriptor)

  envDescriptorCloned.content.files = envDescriptor.content.files.filter(
    envFile => !files.find(file => file === envFile.name)
  )

  if (
    envDescriptorCloned.content.files.length ===
    envDescriptor.content.files.length
  ) {
    throw new Error('No file to delete in keystone')
  }

  return uploadFileInProject(userSession, {
    projectDescriptor,
    fileDescriptor: envDescriptorCloned,
    type: 'env',
    env,
    envDescriptor: envDescriptorCloned,
    versioning: true,
  })
}

module.exports = {
  createEnv,
  removeEnv,
  getEnv,
  // getEnvDescriptor,
  checkoutEnv,
  updateMembersRoles,
  getEnvsDescriptor,
  uploadFileToEnv,
  deleteFilesFromEnv,
  getAdminAndContributorIds,
}
