// const fs = require('fs')
const hash = require('object-hash')

const {
  getPath,
  deleteFiles,
  uploadFile,
  changePathToCurrentUser,
  getFiles,
  getPubkey,
  uploadFileFromOthers,
  getProjects,
  getFileDescriptor,
  getFileFromGaia,
  writeFileOnGaia,
  uploadFileForMembers,
} = require('../file')
// const { createEnv, getEnvsDescriptor } = require('../env')
const { getInvitations } = require('../invitation')

const PROJECTS_STORE = 'projects.json'

const getAllMembers = (projectDescriptor, editorsOnly = false) => {
  const { members } = projectDescriptor.content
  // if (editorsOnly) {
  //   return [...members.admins, ...members.contributors]
  // }
  // return [...members.admins, ...members.contributors, ...members.readers]

  return members.map(m => m.blockstack_id)
}

const addMember = (members, user) => {
  if (members.find(member => member.blockstack_id === user.blockstack_id)) {
    throw new Error('User already in the project')
  }
  return [...members, { blockstack_id: user.blockstack_id, email: user.email }]
}

const getEnvsDescriptor = async (
  userSession,
  { username = null, projectDescriptor = null } = null
) => {
  try {
    const envsDescriptor = await Promise.all(
      projectDescriptor.content.envs.map(async env => {
        const userData = userSession.loadUserData()
        const path = getPath({
          project: projectDescriptor.content.name,
          env: env.name,
          blockstack_id: userData.username,
          type: 'env',
        })
        let envDescriptor
        if (username) {
          envDescriptor = await getFileFromGaia(userSession, path, {
            username,
            decrypt: true,
          })
        } else {
          envDescriptor = await getFileFromGaia(userSession, path)
        }
        if (envDescriptor) {
          // TODO : get last version of the project descriptor among
          return JSON.parse(envDescriptor)
        }
      })
    )
    return envsDescriptor.filter(env => env)
  } catch (error) {
    throw error
  }
}

// create a new project and save it to Gaïa
const createProject = async (userSession, { name, blockstack_id, email }) => {
  const projectContent = {
    name,
    // {name, path}
    envs: [],
    members: [{ blockstack_id, email }],
    createdBy: { blockstack_id, email },
  }

  const projectPath = getPath({
    project: name,
    filename: null,
    blockstack_id,
    type: 'project',
  })

  const projectDescriptor = {
    path: projectPath,
    checksum: hash(JSON.stringify(projectContent)),
    content: projectContent,
    version: 1,
    author: blockstack_id,
    history: {},
  }
  // save project file on Gaïa
  const putFileOptions = { encrypt: true }
  try {
    await writeFileOnGaia(
      userSession,
      projectPath,
      JSON.stringify(projectDescriptor),
      putFileOptions
    )

    // await createEnv(userSession, {
    //   name: 'default',
    //   blockstack_id,
    //   email,
    //   project: name,
    // })

    return projectDescriptor
  } catch (error) {
    throw error
  }
}

const removeProject = async (userSession, { name, force = false }) => {
  // TODO: should remove all files related to the project
  try {
    const projectDescriptor = await getProjectDescriptor(userSession, {
      project: name,
    })
    const envsDescriptor = await getEnvsDescriptor(userSession, {
      projectDescriptor,
    })

    const projects = await getProjects(userSession)
    try {
      // start by removing every files in the project
      const { deletedFiles } = await deleteFiles(userSession, {
        project: name,
        id: projectDescriptor.id,
        envsDescriptor,
      })

      // if a file failed to be deleted, we can't remove the project
      const failedSome = deletedFiles.findIndex(file => {
        return file.deleted === false
      })
      if (failedSome !== -1 && !force) {
        const failures = deletedFiles.filter(file => {
          return file.deleted === false
        })
        throw new Error(`Unable to delete ${failures.join(', ')}`)
      }

      const userData = userSession.loadUserData()

      if (projectDescriptor.content.members.length > 1) {
        // other people are counting on the project
        // if user is admin, set project status as removed
        console.log('TO IMPLEMENT: removing project used by many people')
      } else {
        // remove project descriptor itself
        const path = getPath({
          project: name,
          filename: null,
          blockstack_id: userData.username,
          type: 'project',
        })
        await userSession.deleteFile(path)
      }
      const projectsUpdate = projects.filter(project => project.name !== name)
      await writeFileOnGaia(
        userSession,
        PROJECTS_STORE,
        JSON.stringify(projectsUpdate),
        { encrypt: true }
      )
      return projectsUpdate
    } catch (error) {
      throw error
    }
  } catch (error) {
    throw error
  }
}

const uploadFilesFromProjectOwner = async (userSession, { project }) => {
  const userData = userSession.loadUserData()
  const projects = await getProjects(userSession)
  const { createdBy } = projects.find(x => x.name === project)

  // Get and upload project descriptor
  const pathToProjectDescriptor = getPath({
    project,
    type: 'project',
    blockstack_id: userData.username,
  })

  const projectDescriptor = JSON.parse(
    await getFileFromGaia(userSession, pathToProjectDescriptor, {
      decrypt: true,
      username: createdBy,
    })
  )
  projectDescriptor.path = pathToProjectDescriptor

  projectDescriptor.content.envs = projectDescriptor.content.envs.map(env => ({
    ...env,
    path: changePathToCurrentUser(env.path, userData.username),
  }))

  await uploadFile(userSession, {
    fileDescriptor: projectDescriptor,
    path: pathToProjectDescriptor,
    content: JSON.stringify(projectDescriptor),
  })

  const files = await getFiles(userSession, {
    username: createdBy,
    filesOnly: false,
  })

  return Promise.all(
    files.map(async file => {
      if (!file.error) {
        file = file.descriptor || file
        file.path = changePathToCurrentUser(file.path, userData.username)
        try {
          await uploadFile(userSession, {
            fileDescriptor: file,
            path: file.path,
            content: JSON.stringify(file),
          })
        } catch (err) {
          file.error = err
        }
      }
      return file
    })
  )
}

const copyFilesForMembers = async (userSession, { projectDescriptor }) => {
  const envsDescriptor = await getEnvsDescriptor(userSession, {
    projectDescriptor,
  })

  // get all files
  const members = getAllMembers(projectDescriptor)

  return Promise.all(
    await members.map(async member => {
      // const { blockstack_id } = member
      const blockstackId = member
      const pubkey = await getPubkey(userSession, { blockstackId })

      const files = []

      await Promise.all(
        await envsDescriptor.map(async envDescriptor => {
          files.push(await getFiles(userSession, { envDescriptor }))

          const path = getPath({
            type: 'env',
            env: envDescriptor.content.name,
            blockstack_id: blockstackId,
            project: projectDescriptor.content.name,
          })
          await uploadFile(userSession, {
            path,
            encrypt: pubkey,
            content: JSON.stringify({ ...envDescriptor, path }),
            fileDescriptor: envDescriptor,
          })
        })
      )
      const projectPath = getPath({
        type: 'project',
        blockstack_id: blockstackId,
        project: projectDescriptor.content.name,
      })
      console.log('PROJECTPATH', projectPath)
      await uploadFile(userSession, {
        path: projectPath,
        encrypt: pubkey,
        content: JSON.stringify({ ...projectDescriptor, path: projectPath }),
        fileDescriptor: projectDescriptor,
      })

      const filesWritten = await Promise.all(
        files.map(async file => {
          if (file.fetched) {
            try {
              const path = getPath({
                project: projectDescriptor.content.name,
                filname: file.name,
                blockstack_id: blockstackId,
              })
              await uploadFile(userSession, {
                path,
                encrypt: pubkey,
                content: JSON.stringify(file.descriptor),
                fileDescriptor: file.descriptor,
                versioning: false,
              })
              return {
                file: file.name,
                written: true,
                for: blockstackId,
              }
            } catch (error) {
              return {
                file: file.name,
                written: false,
                for: blockstackId,
                error: error.message,
              }
            }
          }
          return {
            file: file.name,
            written: false,
            for: blockstackId,
            error: file.error,
          }
        })
      )

      return {
        id: blockstackId,
        files: filesWritten,
      }
    })
  )
}

const addMemberToProject = async (
  userSession,
  { projectDescriptor, invitee }
) => {
  const { id, email } = invitee
  // get invitations. We can only add people with invitations open.
  let founByEmail = false
  const invitations = await getInvitations(userSession)

  const member = invitations.find(invite => {
    // find by blockstack_id or email
    if (invite.project === projectDescriptor.content.name) {
      if (invite.blockstack_id === id) {
        return true
      }
      if (invite.email === email) {
        founByEmail = true
        return true
      }
    }
    return false
  })

  if (founByEmail) {
    member.blockstack_id = id
  }

  // keys are plural
  projectDescriptor.content.members = addMember(
    projectDescriptor.content.members,
    member
  )

  return projectDescriptor
}

const addMembersToProject = async (userSession, { project, invitees }) => {
  try {
    let projectDescriptor = await getProjectDescriptor(userSession, { project })

    const projectChecksum = hash(projectDescriptor)

    // TODO if member already in the project

    const userData = userSession.loadUserData()

    const membersAdded = await Promise.all(
      invitees.map(async invitee => {
        console.log('INVITEE', invitee)
        // const path = getPath(project, file.name, userData.username)
        const { id } = invitee
        try {
          // get their public key, if they don't they are not ready
          await getPubkey(userSession, { blockstack_id: id })

          // projectDescriptor.content.members = projectDescriptor.content.members.filter(
          //   m => m.blockstack_id === userData.blockstack_id
          // )

          projectDescriptor = await addMemberToProject(userSession, {
            projectDescriptor,
            invitee,
          })

          return {
            membersAdded: {
              id,
              added: true,
            },
          }
        } catch (error) {
          return {
            id,
            added: false,
            error: error.message,
          }
        }
      })
    )
    await copyFilesForMembers(userSession, { projectDescriptor })
    const newProjectChecksum = hash(projectDescriptor)

    if (newProjectChecksum !== projectChecksum) {
      // project descriptor has been updated

      await uploadFile(userSession, {
        projectDescriptor,
        project: projectDescriptor.name,
        fileDescriptor: projectDescriptor,
        file: project,
        versioning: true,
        type: 'project',
        encryptForOthers: true,
        members: getAllMembers(projectDescriptor),
      })
    }

    return { membersAdded, projectDescriptor }
  } catch (error) {
    throw error
  }
}

const updateProjectsStatus = async userSession => {
  const projectsFiles = await getProjects(userSession)

  const userData = userSession.loadUserData()
  const projects = await Promise.all(
    projectsFiles.map(async project => {
      if (!project.pendingInvite) {
        return project
      }
      try {
        const path = getPath({
          blockstack_id: userData.username,
          project: project.name,
          type: 'project',
        })
        const file = await getFileFromGaia(userSession, path, {
          username: project.createdBy,
          decrypt: true,
        })
        if (!file) {
          throw new Error(
            `The project file ${project.name} is not available in ${project.createdBy} workspace`
          )
        }

        return { ...project, pendingInvite: false }
      } catch (err) {
        return { ...project, pendingInvite: true }
      }
    })
  )
  return projects
}

const updateProjectsFromStatuses = async (
  userSession,
  { projects, projectsStatuses }
) => {
  const projectsChecksum = hash(projects)
  projects = await Promise.all(
    projects.map(async project => {
      const projectUpdate = projectsStatuses.find(
        ps => ps.project.name === project.name
      )
      if (
        projectUpdate &&
        projectUpdate.invite === 'fulfilled' &&
        projectUpdate.projectDescriptor
      ) {
        const members = getAllMembers(projectUpdate.projectDescriptor)

        await uploadFileFromOthers(userSession, {
          projectDescriptor: projectUpdate.projectDescriptor,
        })
        return {
          ...project,
          pendingInvite: false,
          at: members,
        }
      }
      return project
    })
  )
  const projectsUpdateChecksum = hash(projects)

  if (projectsChecksum !== projectsUpdateChecksum) {
    // is there any changes? yes? save the file online.
    await writeFileOnGaia(
      userSession,
      PROJECTS_STORE,
      JSON.stringify(projects),
      {
        encrypt: true,
      }
    )
  }

  return projects
}

const getAllMembersNotSelf = (
  userSession,
  { projectDescriptor, editorsOnly = false }
) => {
  const userData = userSession.loadUserData()
  const members = getAllMembers(projectDescriptor, editorsOnly)
  return members.filter(member => member.blockstack_id !== userData.username)
}

const uploadFileInProject = (
  userSession,
  { projectDescriptor, fileDescriptor, type, env, versioning, members }
) => {
  // TODO: save pubkey in order to save requests

  if (!members) {
    members = getAllMembers(projectDescriptor)
  }

  return uploadFileForMembers(userSession, {
    project: projectDescriptor.name,
    env,
    type,
    fileDescriptor,
    versioning,
    members,
  })
}

const getProjectDescriptor = async (userSession, { project }) => {
  try {
    const userData = userSession.loadUserData()
    // if (!project) {
    //   project = JSON.parse(fs.readFileSync('.ksconfig')).project
    // }
    const path = getPath({
      project,
      blockstack_id: userData.username,
      type: 'project',
    })
    const descriptor = await getFileFromGaia(userSession, path)
    if (descriptor) {
      // TODO : get last version of the project descriptor among
      return JSON.parse(descriptor)
    }
    throw new Error(`Project ${project} not found`)
  } catch (error) {
    throw error
  }
}

const syncFilesWithAfterBeingInvited = async (
  userSession,
  { project, invitedBy }
) => {
  const userData = userSession.loadUserData()

  // Retrieve project descriptor from invited by.
  const projectDescriptor = await getFileDescriptor(userSession, {
    author: invitedBy,
    type: 'project',
    project,
  })

  const projectPath = getPath({
    type: 'project',
    project,
    blockstack_id: userData.username,
  })

  await uploadFile(userSession, {
    fileDescriptor: { ...projectDescriptor, path: projectPath },
  })

  const envsDescriptor = await getEnvsDescriptor(userSession, {
    username: invitedBy,
    projectDescriptor,
  })
  await Promise.all(
    envsDescriptor.map(async envDescriptor => {
      const envPath = await getPath({
        env: envDescriptor.content.name,
        type: 'env',
        project,
        blockstack_id: userData.username,
      })
      await uploadFile(userSession, {
        fileDescriptor: { ...envDescriptor, path: envPath },
      })
    })
  )
}

module.exports = {
  createProject,
  removeProject,
  uploadFilesFromProjectOwner,
  addMembersToProject,
  updateProjectsStatus,
  // getProjectName,
  updateProjectsFromStatuses,
  getAllMembersNotSelf,
  // uploadFileForOthers,
  uploadFileInProject,
  getProjectDescriptor,
  getEnvsDescriptor,
  syncFilesWithAfterBeingInvited,
  getAllMembers,
}
