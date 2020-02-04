const isUUID = require('uuid-validate')
const { readFileFromGaia, writeFileToGaia } = require('../file/gaia')
const { getPath } = require('../descriptor-path')
// const { getInvitations, isMemberInvited } = require('../invitation')
const KeystoneError = require('../error')
const { ROLES, PROJECTS_STORE } = require('../constants')
const {
  updateDescriptorForMembers,
  updateDescriptor,
  getDescriptor,
} = require('../descriptor')

const createProjectsStore = async userSession => {
  return writeFileToGaia(userSession, {
    path: PROJECTS_STORE,
    content: JSON.stringify([]),
  })
}

const getProjects = async userSession => {
  const projectsFile = await readFileFromGaia(userSession, {
    path: PROJECTS_STORE,
  })
  if (projectsFile) return projectsFile

  // if null, it means projects.json hasn't been found on Gaia. We need to create it.
  createProjectsStore(userSession)
  return []
}

const findProjectByUUID = (projects, name) => {
  try {
    const [_, uuid] = getNameAndUUID(name)
    return projects.find(p => p.name.indexOf(`/${uuid}`) > 0)
  } catch (err) {
    return undefined
  }
}
// is there a project with the same name?
const findProjectByName = (projects, name) =>
  projects.filter(p => p.name.indexOf(`${name}/`) === 0)

const createProject = async (
  userSession,
  { name, members, pendingInvite = false }
) => {
  const { username } = userSession.loadUserData()

  // TODO: validate members format
  const project = {
    name,
    members,
    createdBy: username,
    pendingInvite,
    env: ['default'],
  }

  const projects = await getProjects(userSession)

  const duplicate = findProjectByName(projects, name)

  if (duplicate.length > 0) throw new Error('The project is already created')

  projects.push(project)

  // update Projects store
  await writeFileToGaia(userSession, {
    path: PROJECTS_STORE,
    content: JSON.stringify(projects),
  })

  const descriptorPath = getPath({
    type: 'project',
    project: name,
    blockstackId: username,
  })

  const membersDescriptor = {
    content: {
      [ROLES.ADMINS]: [{ blockstack_id: username }],
      [ROLES.CONTRIBUTORS]: [],
      [ROLES.READERS]: [],
    },
  }

  // create project descriptor and update user storage space for every members
  const projectDescriptor = await updateDescriptorForMembers(userSession, {
    descriptorPath,
    project: name,
    name,
    type: 'project',
    membersDescriptor,
    content: { env: ['default'] },
  })

  return projectDescriptor
}

const syncProjectsStatus = async userSession => {
  const projectsFiles = await getProjects(userSession)

  const projects = await Promise.all(
    projectsFiles.map(async project => {
      console.log(project)
      //   if (!project.pendingInvite) {
      //     return project
      //   }
      //   try {
      //     const path = getPath({
      //       blockstack_id: userData.username,
      //       project: project.name,
      //       type: 'project',
      //     })
      //     const file = await getFileFromGaia(userSession, path, {
      //       username: project.createdBy,
      //       decrypt: true,
      //     })
      //     if (!file) {
      //       throw new Error(
      //         `The project file ${project.name} is not available in ${project.createdBy} workspace`
      //       )
      //     }

      //     return { ...project, pendingInvite: false }
      //   } catch (err) {
      //     return { ...project, pendingInvite: true }
      //   }
    })
  )
  return projects
}

// const addMember = async (descriptor, user) => {
//   const allMembers = await getMembers(userSession, { project, env })

//   if (allMembers.find(member => member === user.blockstackId)) {

//   }

//   return {
//     ...descriptor.content.members,
//     [`${user.role}s`]: [
//       ...descriptor.content.members[`${user.role}s`],
//       { email: user.email, blockstack_id: user.blockstackId },
//     ],
//   }
// }

const addEnvToProject = (userSession, { projectDescriptor, env }) => {
  const newProjectDescriptor = { ...projectDescriptor }

  newProjectDescriptor.content.env.push(env)

  return updateDescriptor(userSession, {
    descriptorPath: projectDescriptor.path,
    project: projectDescriptor.name,
    type: 'project',
    content: newProjectDescriptor.content,
    name: projectDescriptor.name,
  })
}

const removeEnvFromProject = async (userSession, { project, env }) => {
  const { username } = userSession.loadUserData()

  const projectDescriptor = await getDescriptor(userSession, {
    project,
    type: 'project',
    blockstackId: username,
  })

  const newProjectDescriptor = { ...projectDescriptor }

  newProjectDescriptor.content.env = newProjectDescriptor.content.env.filter(
    e => e !== env
  )

  return updateDescriptor(userSession, {
    descriptorPath: projectDescriptor.path,
    project: projectDescriptor.name,
    type: 'project',
    content: newProjectDescriptor.content,
    name: projectDescriptor.name,
  })
}

const getNameAndUUID = projectFullname => {
  try {
    const projectParts = projectFullname.split('/')
    const name = projectParts[0]
    const uuid = projectParts[1]
    // UUID should be in version 4
    if (!isUUID(uuid, 4)) throw new Error('UUID missing')
    return [name, uuid]
  } catch (err) {
    throw new KeystoneError(
      'InvalidProjectName',
      'Invalid project name',
      projectFullname
    )
  }
}

module.exports = {
  syncProjectsStatus,
  getProjects,
  createProject,
  findProjectByName,
  findProjectByUUID,
  addEnvToProject,
  removeEnvFromProject,
  getNameAndUUID,
}
