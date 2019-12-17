// const fs = require('fs')
const daffy = require('daffy')
const filenameRegex = require('filename-regex')
const hash = require('object-hash')
const debug = require('debug')('keystone:core:file')
const NodeCache = require('node-cache')

const fileCache = new NodeCache()
const PUBKEY = 'public.key'
const PROJECTS_STORE = 'projects.json'

const getFileFromGaia = async (userSession, path, opts) => {
  const userData = userSession.loadUserData()
  debug('')
  const cacheKey = `${(opts && opts.username) || userData.username}/${path}`
  const file = fileCache.get(cacheKey)
  if (file) {
    debug(
      `Get file from cache ${path} from ${
        opts && opts.username ? opts.username : 'self'
      }`
    )
    return file
  }
  debug(
    `Get file from gaia ${path} from ${
      opts && opts.username ? opts.username : 'self'
    }`
  )
  const fetchedFile = await userSession.getFile(path, opts)
  fileCache.set(cacheKey, fetchedFile)
  return fetchedFile
}

const writeFileOnGaia = async (userSession, path, content, opts) => {
  const userData = userSession.loadUserData()

  debug('')
  debug(
    `Write file on gaia to ${path} with encrypt ${
      opts ? opts.encrypt : 'no encrypt'
    }`
  )
  const cacheKey = `${(opts && opts.username) || userData.username}/${path}`
  await userSession.putFile(path, content, opts)

  fileCache.set(cacheKey, content)

  return content
}

const getEnvPath = ({ env, project, blockstack_id }) => {
  if (!project) {
    throw new Error('Need project!')
  }
  if (!env) {
    throw new Error('Need env!')
  }

  return `${project}/${env}/${blockstack_id}.json`
}

const getFilePath = ({ filename, blockstack_id, env, project }) => {
  if (!project) {
    throw new Error('Need project!')
  }
  if (!env) {
    throw new Error('Need env!')
  }
  if (!filename) {
    throw new Error('Need filename!')
  }

  return `${project}/${env}/${filename}/${blockstack_id}.json`
}

const getProjectPath = ({ blockstack_id, project }) => {
  if (!project) {
    throw new Error('Need project!')
  }

  return `${project}/${blockstack_id}.json`
}

const getPath = ({ env, project, blockstack_id, type, filename }) => {
  let path

  switch (type) {
    case 'project':
      path = getProjectPath({
        blockstack_id,
        project,
      })
      break

    case 'env':
      path = getEnvPath({
        blockstack_id,
        env,
        project,
      })
      break

    case 'file':
      path = getFilePath({
        blockstack_id,
        env,
        project,
        filename,
      })
      break

    default:
      throw new Error(`Upload file with unknown type: ${type}`)
  }

  return path
}

const extractFileFromPath = file => {
  const filename = file.match(filenameRegex())
  return filename[0]
}

// filename can be the project not only files in the project.
const uploadFile = async (
  userSession,
  {
    path,
    fileDescriptor,
    content,
    file,
    encrypt = true,
    versioning = false,
    type = 'file',
    // encryptForOthers = true,
    members,
    project,
    env,
  }
) => {
  if (versioning) {
    // const projectDescriptor = await getProjectDescriptor(userSession, {})
    const filename = extractFileFromPath(
      file || fileDescriptor.name || fileDescriptor.content.name
    )
    // start by syncing the file
    const lastVersion = await syncFile(userSession, {
      fileDescriptor,
      filename,
      type,
      members,
      project,
      env,
    })

    await writeFileOnGaia(userSession, path, JSON.stringify(lastVersion), {
      encrypt,
    })
    return lastVersion
  }

  const result = await writeFileOnGaia(
    userSession,
    path || fileDescriptor.path,
    content || JSON.stringify(fileDescriptor),
    {
      encrypt,
    }
  )

  return result
}

const getProjects = async userSession => {
  try {
    const projectsFile = await getFileFromGaia(userSession, PROJECTS_STORE)
    if (projectsFile) {
      return JSON.parse(projectsFile)
    }
    throw new Error(`Projects not found`)
  } catch (error) {
    throw error
  }
}

// if files is empty, delete every files from the project
const deleteFiles = async (
  userSession,
  { files = [], envs = [], members, envsDescriptor }
) => {
  try {
    // const envsDescriptor = await getEnvsDescriptor(
    //   userSession,
    //   {}
    // )
    const envsChecksum = envsDescriptor.map(env => ({
      name: env.name,
      checksum: env.checksum,
    }))
    files = await getFiles(userSession, { envs, files, envsDescriptor })

    // const projectDescriptor = await getProjectDescriptor(userSession, {})
    // const members = projectDescriptor.content.members.length
    if (members > 1) {
      // other people are using the files
      // if user is admin, set files status as removed
      console.log('TO IMPLEMENT: removing project used by many people')
    } else {
      // we can safely remove everything as the user
      // is the only one using the project

      // // start by removing files from the project
      const deletedFiles = await Promise.all(
        files.map(async file => {
          try {
            let fileEnvDescriptor = envsDescriptor.find(
              e => e.content.name === file.env
            )
            await userSession.deleteFile(file.descriptor.path)

            fileEnvDescriptor = removeFileFromEnv({
              envDescriptor: fileEnvDescriptor,
              file: file.name,
            })
            envsDescriptor[
              envsDescriptor.findIndex(e => e.name === file.env)
            ] = fileEnvDescriptor
            return {
              name: file.name,
              path: file.descriptor.path,
              deleted: true,
            }
          } catch (error) {
            console.log(error)
            return {
              name: file.name,
              path: file.descriptor.path,
              deleted: false,
              error: error.message,
            }
          }
        })
      )

      await Promise.all(
        envsDescriptor.map(async envDescriptor => {
          const newEnvChecksum = hash(envDescriptor)

          if (
            newEnvChecksum !==
            envsChecksum[
              envsChecksum.findIndex(env => env.name === envDescriptor.name)
            ].checksum
          ) {
            // project descriptor has been updated
            await uploadFile(userSession, {
              fileDescriptor: envDescriptor,
              file: envDescriptor.content.name,
              versioning: true,
              type: 'env',
            })
          }
        })
      )

      return { deletedFiles, envsDescriptor }
    }
  } catch (error) {
    throw error
  }
}

const getFiles = async (
  userSession,
  { files = [], username, filesOnly = true, envDescriptor, project }
) => {
  try {
    const filesToDisplay = envDescriptor.content.files.filter(file =>
      files.find(fileToCat => fileToCat === file.name)
    )
    // let envsDescriptor = await getEnvsDescriptor(userSession, {
    //   username,
    // })

    // const filteredEnvsDescriptor = envsDescriptor.filter(
    //   env => envs.length === 0 || envs.includes(env.content.name)
    // )

    // let projectFiles = filteredEnvsDescriptor.reduce((envsFiles, descriptor) => {
    //   envsFiles.push(
    //     ...descriptor.content.files.map(file => ({
    //       ...file,
    //       env: descriptor.content.name,
    //     }))
    //   )

    //   return envsFiles
    // }, [])

    // projectFiles = projectFiles.filter(
    //   file => files.includes(file.name) || files.length === 0
    // )

    const userData = userSession.loadUserData()

    const fetchedFiles = await Promise.all(
      filesToDisplay.map(async file => {
        // const path = getPath(project, file.name, userData.username)
        try {
          const fileDescriptor = await getFileDescriptor(userSession, {
            file: file.name,
            author: username || userData.username,
            env: envDescriptor.content.name,
            blockstack_id: userData.username,
            type: 'file',
            project,
          })
          console.log('TCL: fileDescriptor', fileDescriptor)
          // const lastVersion = await syncFile(userSession, {
          //   projectDescriptor,
          //   fileDescriptor,
          //   filename: file.name,
          //   preflight: true,
          // })

          return {
            name: file.name,
            fetched: true,
            descriptor: fileDescriptor,
            env: file.env,
          }
        } catch (error) {
          console.error(error)
          return {
            name: file.name,
            fetched: false,
            content: null,
            error: error.message,
          }
        }
      })
    )

    // if (!filesOnly) {
    //   fetchedFiles.push(...envsDescriptor)
    // }

    return fetchedFiles
  } catch (error) {
    throw error
  }
}

/**
 * Return file descriptor
 * @param {*} param1.author User you want to get your file from
 */
const getFileDescriptor = async (
  userSession,
  { file: filename, author, type, env, project }
) => {
  const userData = userSession.loadUserData()

  const path = getPath({
    filename,
    blockstack_id: userData.username,
    type,
    env,
    project,
  })

  const descriptor = await getFileFromGaia(userSession, path, {
    username: author || userData.usrname,
    decrypt: true,
  })
  if (descriptor) {
    return JSON.parse(descriptor)
  }
  throw new Error(`Unable to fetch file from storage: ${path}`)
}

const newFileDescriptor = ({
  filename,
  project,
  content,
  author,
  type = 'file',
  env,
}) => {
  return {
    path: getPath({ project, filename, blockstack_id: author, type, env }),
    name: filename,
    content,
    checksum: hash(content),
    history: [],
    author,
    version: 0,
  }
}

const getPubkey = async (userSession, { blockstack_id }) => {
  try {
    const pubkeyFile = await getFileFromGaia(userSession, PUBKEY, {
      username: blockstack_id,
      decrypt: false,
    })
    if (pubkeyFile) {
      return pubkeyFile
    }
    throw new Error(
      `Keystone public application key not found on ${blockstack_id}`
    )
  } catch (error) {
    console.log()
    throw error
  }
}

// update projects.json
const updateWorkspace = async (
  userSession,
  { name, blockstack_id, invitation = false, at = null }
) => {
  const project = {
    name,
    at,
    createdBy: blockstack_id,
    pendingInvite: !!invitation,
  }

  // retrieve projects.json
  const projects = await getProjects(userSession)

  projects.push(project)

  try {
    await writeFileOnGaia(
      userSession,
      PROJECTS_STORE,
      JSON.stringify(projects),
      {
        encrypt: true,
      }
    )
  } catch (error) {
    console.log("Couldn't save workspace", error)
  }

  return projects
}

const getFileFromOther = async (
  userSession,
  { username, origin, project, filename, type, env }
) => {
  // path is made of 3 parts: project, file, username used for encryption
  // my-project/my-file.json/my-blockstack.id.blockstack
  const path = getPath({
    project,
    filename,
    blockstack_id: username,
    type,
    env,
  })
  try {
    const fileDescriptor = await getFileFromGaia(userSession, path, {
      username: origin,
      decrypt: true,
    })

    if (fileDescriptor) {
      return {
        id: username,
        fileDescriptor: JSON.parse(fileDescriptor),
        fetched: true,
      }
    }
    throw new Error(`${filename} not found on ${username}`)
  } catch (error) {
    return {
      fetched: false,
      error,
      fileDescriptor: null,
      id: username,
    }
  }
}

const getFileFromEveryone = async (
  userSession,
  { fileDescriptor, members, type, project, env }
) => {
  const userData = userSession.loadUserData()
  const filesToFetch = members
    .filter(member => member !== userData.username)
    .map(async member =>
      getFileFromOther(userSession, {
        origin: member,
        project,
        filename: fileDescriptor.name || fileDescriptor.content.name,
        type,
        env,
        username: userData.username,
      })
    )

  const filesFetched = await Promise.all([...filesToFetch])

  const files = filesFetched.reduce((filesAcc, fileFetched) => {
    if (fileFetched.fetched) {
      return [...filesAcc, fileFetched]
    }
    return filesAcc
  }, [])

  // find greatest version
  let currentUserVersion
  let version = {
    conflict: false,
    stable: null,
    versions: [],
  }

  if (files.length === 0) {
    try {
      await uploadFileForMembers(userSession, {
        project,
        env,
        type,
        fileDescriptor,
        versioning: false,
        members,
      })
    } catch (err) {}
    return {
      files,
      version: {
        conflict: false,
        stable: fileDescriptor,
        isUserVersion: true,
      },
    }
  }

  try {
    currentUserVersion = await getFileDescriptor(userSession, {
      blockstack_id: userData.username,
      file: fileDescriptor.name,
      project,
      author: userData.username,
      type,
      env,
    })
  } catch (err) {
    console.error(err)
  }

  version = getGreatestVersion(
    files.map(file => file.fileDescriptor),
    currentUserVersion
  )
  if (!version.conflict && !version.isUserVersion) {
    await uploadFileForMembers(userSession, {
      project,
      env,
      type,
      fileDescriptor,
      versioning: false,
      members,
    })
  }
  return { files, version }
}

const hasConflicts = (versions, index) => {
  if (versions[index].length === 1) return false
  // same version with different content
  if (
    versions[index].find(
      version => version.checksum !== versions[index][0].checksum
    )
  ) {
    return true
  }
  return false
}

const getGreatestStableVersion = (versions, index) => {
  if (hasConflicts(versions, index)) {
    const prev = index - 1
    if (prev >= 1) {
      return getGreatestStableVersion(versions, prev)
    }
    throw new Error('No stable version found.')
  }
  return versions[index][0]
}

const uploadFileFromOthers = async (userSession, { projectDescriptor }) => {
  const userData = userSession.loadUserData()
  const path = getPath({
    project: projectDescriptor.content.name,
    blockstack_id: userData.username,
    type: 'project',
  })

  // upload project descriptor in own store
  await uploadFile(userSession, {
    path,
    projectDescriptor,
    fileDescriptor: projectDescriptor,
    file: projectDescriptor.content,
    versioning: false,
    type: 'project',
  })

  const envsDescriptor = await Promise.all(
    projectDescriptor.content.envs.map(async env => {
      const { name } = env
      const envPath = getPath({
        project: projectDescriptor.content.name,
        type: 'env',
        env: name,
        blockstack_id: userData.username,
      })
      const envDescriptor = await getFileFromGaia(userSession, envPath, {
        username: { blockstack_id: projectDescriptor.author },
      })
      console.log('envDescriptor', envDescriptor)
      return envDescriptor
    })
  )
  console.log(envsDescriptor)
}

const incrementVersion = ({
  fileDescriptor,
  author,
  lastDescriptor = null,
  type,
}) => {
  const { content } = fileDescriptor
  const newChecksum = hash(content)

  if (lastDescriptor) {
    // same content, no need to update
    if (newChecksum === lastDescriptor.checksum) {
      // we avoid throwing an error for project files
      // as it would happens everytime a user push files.
      if (type !== 'project' && type !== 'env') {
        throw new Error(
          'A version of this file with the same content already exists.'
        )
      }
    }
    const newEntry = {
      version: lastDescriptor.version,
      checksum: lastDescriptor.checksum,
      content: daffy.createPatch(
        JSON.stringify(content),
        JSON.stringify(lastDescriptor.content)
      ),
      sourcePatch: newChecksum,
      author: lastDescriptor.author,
    }

    const history =
      lastDescriptor.history && lastDescriptor.history.length > 0
        ? lastDescriptor.history
        : []

    return {
      ...fileDescriptor,
      checksum: newChecksum,
      version: lastDescriptor.version + 1,
      history: [...history, newEntry],
      author,
    }
  }

  return {
    ...fileDescriptor,
    checksum: newChecksum,
    version: 1,
    history: [],
    author,
  }
}

const getGreatestVersion = (fileFromMembers, userFile = { version: -1 }) => {
  // group by version
  const versions = fileFromMembers.reduce((grouped, file) => {
    const newGrouped = grouped
    newGrouped[file.version.toString()] = [
      ...(newGrouped[file.version.toString()] || []),
      file,
    ]
    return newGrouped
  }, {})

  // add user version to the mix.
  versions[userFile.version.toString()] = [
    ...(versions[userFile.version.toString()] || []),
    userFile,
  ]

  const greatestVersionNumber = Math.max(
    ...Object.keys(versions).map(x => parseInt(x))
  )

  try {
    if (hasConflicts(versions, greatestVersionNumber)) {
      return {
        conflict: true,
        versions: versions[greatestVersionNumber],
        stable: getGreatestStableVersion(),
      }
    }
  } catch (error) {
    return {
      conflict: true,
      versions: versions[greatestVersionNumber],
      stable: false,
    }
  }

  return {
    conflict: false,
    stable: versions[greatestVersionNumber][0],
    versions: null,
    isUserVersion:
      userFile.checksum === versions[greatestVersionNumber][0].checksum,
  }
}

// search for newer version in members workspaces
// if no conflict, increment the version and the history
// preflight = make everything like real but don't increment the version
const syncFile = async (
  userSession,
  { project, env, members, fileDescriptor, type, preflight = false }
) => {
  const userData = userSession.loadUserData()
  let files = []
  let version = { conflict: false, stable: null, version: [] }
  try {
    const fileFromOther = await getFileFromEveryone(userSession, {
      fileDescriptor,
      members,
      type,
      project,
      env,
    })
    files = fileFromOther.files
    version = fileFromOther.version
  } catch (err) {
    console.log(err)
  } finally {
    if (version.conflict) {
      // handle conflict there
      // if 2 versions with different checksum, try to merge.
      // if failed ask the user to resolve the conflict (admin or contributor)
      console.log('CONFLICTS')
      return version
    }
    if (preflight) {
      return version
    }

    if (
      files.length === 0 ||
      version.stable.checksum === fileDescriptor.checksum
    ) {
      return version.stable
    }
    // the current user is the only owner

    // we have known history
    if (version.stable) {
      return incrementVersion({
        fileDescriptor,
        author: userData.username,
        lastDescriptor: version.stable,
        type,
      })
    }
    return incrementVersion({
      fileDescriptor,
      author: userData.username,
      type,
    })
    // }
  }
}

// const extractFileFromPath = file => {
//   const filename = file.match(filenameRegex())
//   return filename[0]
// }

// const getProjectDescriptor = async (userSession, { project }) => {
//   try {
//     const userData = userSession.loadUserData()
//     if (!project) {
//       project = JSON.parse(fs.readFileSync('.ksconfig')).project
//     }
//     const path = getPath({
//       project,
//       blockstack_id: userData.username,
//       type: 'project',
//     })
//     const descriptor = await userSession.getFile(path)
//     if (descriptor) {
//       // TODO : get last version of the project descriptor among
//       return JSON.parse(descriptor)
//     }
//     throw new Error(`Project ${project} not found`)
//   } catch (error) {
//     throw error
//   }
// }

const removeFileFromEnv = ({ envDescriptor, file }) => {
  const filename = extractFileFromPath(file)

  // remove the file from the project descriptor
  envDescriptor.content.files = envDescriptor.content.files.filter(
    f => f.name !== filename
  )

  console.log('removeFileFromEnv', envDescriptor)
  return {
    ...envDescriptor,
    content: {
      ...envDescriptor.content,
      files: envDescriptor.content.files,
    },
  }
}

async function uploadFileForMembers(
  userSession,
  { project, env, type, fileDescriptor, versioning, members }
) {
  // console.log('uploadedFile', result)
  return Promise.all(
    members.map(async blockstack_id => {
      try {
        const pubkey = await getPubkey(userSession, {
          blockstack_id,
        })
        // console.log({ projectDescriptor, file, filename, type })
        const memberPath = getPath({
          blockstack_id,
          project,
          env,
          type,
          filename: fileDescriptor.name,
        })
        await uploadFile(userSession, {
          project,
          path: memberPath,
          content: JSON.stringify({ path: memberPath, ...fileDescriptor }),
          fileDescriptor,
          encrypt: pubkey,
          members,
          env,
          type,
          versioning,
        })
        return {
          id: blockstack_id,
          uploaded: true,
          file: fileDescriptor.name || fileDescriptor.content.name,
        }
      } catch (error) {
        console.error(error)
        return {
          id: blockstack_id,
          uploaded: false,
          file: fileDescriptor.name || fileDescriptor.content.name,
          error: error.message,
        }
      }
    })
  )
}

module.exports = {
  uploadFile,
  newFileDescriptor,
  deleteFiles,
  getFiles,
  getPubkey,
  updateWorkspace,
  getProjectPath,
  getEnvPath,
  getFilePath,
  uploadFileFromOthers,
  // getProjectDescriptor,
  removeFileFromEnv,
  getFileDescriptor,
  getPath,
  getProjects,
  extractFileFromPath,
  getFileFromEveryone,
  getFileFromGaia,
  writeFileOnGaia,
  uploadFileForMembers,
}
