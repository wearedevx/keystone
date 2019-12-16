const axios = require('axios')
const isEmail = require('is-email')
const KeystoneError = require('../error')
const {
  KEYSTONE_MAIL,
  INVITATIONS_STORE,
  PROJECTS_STORE,
} = require('../constants')
const { readFileFromGaia, writeFileToGaia } = require('../file/gaia')
const {
  getProjects,
  findProjectByUUID,
  getNameAndUUID,
} = require('../projects')
const { addMember } = require('../member')

const createInvitationsFile = async userSession => {
  try {
    await writeFileToGaia(userSession, {
      path: INVITATIONS_STORE,
      content: JSON.stringify([]),
    })
    return []
  } catch (error) {
    throw error
  }
}

const getInvitations = async userSession => {
  try {
    const invitationsFile = await readFileFromGaia(userSession, {
      path: INVITATIONS_STORE,
    })
    if (invitationsFile) {
      return invitationsFile
    }
    throw new Error(`Invitations file not found`)
  } catch (error) {
    console.log(error)
    // if invitations not found, create the file
    const invitationsFile = await createInvitationsFile(userSession)
    return invitationsFile
  }
}

const updateInvitations = async (userSession, { invitations }) => {
  try {
    await writeFileToGaia(userSession, {
      path: INVITATIONS_STORE,
      content: JSON.stringify(invitations),
    })
    return []
  } catch (error) {
    throw error
  }
}

const addInvitee = async (userSession, { project, email, role }) => {
  // check if invitations.json exists
  try {
    let invitations = await getInvitations(userSession)
    const { username } = userSession.loadUserData()

    // does the invitationn already exists?
    const alreadyInvited = invitations.find(invitation => {
      return invitation.email === email && invitation.project === project
    })

    if (alreadyInvited) {
      throw new Error(`${email} has already been invited to ${project}`)
    }

    invitations = [
      ...invitations,
      {
        email,
        role,
        project,
        blockstack_id: username,
      },
    ]

    return invitations
  } catch (error) {
    throw error
  }
}

const inviteMember = async (
  userSession,
  { from, project, emails, role = 'reader' }
) => {
  const userData = userSession.loadUserData()
  const [projectName, projectUUID] = getNameAndUUID(project)
  const emailsSent = await Promise.all(
    emails.map(async email => {
      // loosely check if format is an email
      if (isEmail(email)) {
        try {
          const invitations = await addInvitee(userSession, {
            project,
            email,
            role,
          })

          await axios({
            method: 'post',
            url: KEYSTONE_MAIL,
            data: {
              request: 'invite',
              from,
              id: userData.username,
              project: projectName,
              role,
              email,
              uuid: projectUUID,
            },
          })

          await updateInvitations(userSession, { invitations })

          return {
            email,
            sent: true,
          }
        } catch (error) {
          return {
            email,
            sent: false,
            error: error.message,
          }
        }
      } else {
        return {
          email,
          sent: false,
          error: `Email address ${email} is invalid`,
        }
      }
    })
  )

  return emailsSent
}

const isMemberInvited = (invitations, project, { blockstack_id, email }) => {
  let foundByEmail

  const member = invitations.find(invitation => {
    // find by blockstack_id or email
    if (invitation.project === project) {
      if (invitation.blockstack_id === blockstack_id) {
        return true
      }
      if (invitation.email === email) {
        foundByEmail = true
        return true
      }
    }
    return false
  })

  if (foundByEmail) {
    member.blockstack_id = blockstack_id
  }

  return member
}

const deleteInvites = async (userSession, { project, emails }) => {
  // check if invitations.json exists
  try {
    let invitations = await getInvitations(userSession)
    let deleted = []

    invitations = invitations.reduce((invites, invite) => {
      const foundInvite = emails.find(email => {
        return email === invite.email && project === invite.project
      })

      if (foundInvite) {
        deleted = [
          ...deleted,
          invitations.find(invitation => {
            return invitation.email === foundInvite && project
          }),
        ]
      }

      if (!foundInvite) {
        return [...invites, invite]
      }
      return invites
    }, [])

    await updateInvitations(userSession, { invitations })
    return deleted
  } catch (error) {
    throw error
  }
}

const acceptInvite = async (
  userSession,
  { name, from, blockstackId, userEmail }
) => {
  const projects = await getProjects(userSession)
  const { username } = userSession.loadUserData()
  const projectFound = findProjectByUUID(projects, name)

  // the project is not in the user projects list
  if (projectFound) {
    throw new KeystoneError(
      'ProjectNameExists',
      'You already have this project in your account',
      name
    )
  } else {
    const project = {
      name,
      createdBy: blockstackId,
      pendingInvite: true,
      env: [],
    }
    projects.push(project)
    try {
      await writeFileToGaia(userSession, {
        path: PROJECTS_STORE,
        content: JSON.stringify(projects),
      })

      const [humanName, uuid] = getNameAndUUID(name)

      await axios({
        method: 'post',
        url: KEYSTONE_MAIL,
        data: {
          request: 'accept',
          from: userEmail,
          to: from,
          id: username,
          project: humanName,
          uuid,
        },
      })

      return projects
    } catch (error) {
      console.log('TCL: error', error)
    }
  }
}

const addMemberToProject = async (userSession, { project, invitee }) => {
  const { role } = invitee
  // get invitations. We can only add people with invitations open.
  const invitations = await getInvitations(userSession)

  const member = isMemberInvited(invitations, project, invitee)
  if (!member) {
    throw new KeystoneError(
      ERROR_CODES.InvitationFailed,
      'User has not been invited to the project',
      invitee
    )
  }

  return addMember(userSession, {
    project,
    member: member.blockstack_id,
    role: `${role}s`,
  })
}

module.exports = {
  inviteMember,
  getInvitations,
  isMemberInvited,
  deleteInvites,
  acceptInvite,
  addMemberToProject,
}
