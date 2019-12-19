import React, { useState } from 'react'
import useUser from '../hooks/useUser'
import queryString from 'query-string'
import ErrorCard from '../components/cards/error'
import BaseCard from '../components/cards/base'
import Button from '../components/button'
import { getNameAndUUID } from '@keystone/core/lib/projects'
import { acceptInvite } from '@keystone/core/lib/invitation'
import WithLoggin from '../components/withLoggin'

const TitlePromptInvite = ({ project }) => (
  <>
    <span
      role="img"
      aria-label="The back of an envelope, as used to send a letter or card in the mail."
      className="mr-2"
    >
      ✉️
    </span>
    You are invited to project <strong>{project}</strong>.
  </>
)

const join = async (
  userSession,
  { name, from, blockstackId, userEmail, onError, onDone, onSuccess }
) => {
  try {
    const projects = await acceptInvite(userSession, {
      name,
      from,
      blockstackId,
      userEmail,
    })
    onSuccess(projects)
    return projects
  } catch (error) {
    onError(error)
  } finally {
    onDone()
  }
}

const PromptInvite = ({
  project,
  uuid,
  adminUserEmail,
  blockstackId,
  userEmail,
}) => {
  const { userSession } = useUser()
  const [joining, setJoining] = useState(false)
  const [error, setError] = useState(false)
  const [success, setSuccess] = useState(false)

  return (
    <>
      <BaseCard title={<TitlePromptInvite project={project} />}>
        {joining && <p>Updating your projects list, please wait...</p>}

        {!joining && (
          <p>
            This invite is sent by <strong>{adminUserEmail}</strong>. Click join
            to join the project or ignore if you don't know the sender.
          </p>
        )}

        <p className="italic text-red-400 mt-6 text-xs uppercase">
          Project id: {uuid}
        </p>
      </BaseCard>

      {error && (
        <p className="text-red-600 font-bold mt-4">
          <span
            role="img"
            aria-label="A triangle with an exclamation mark inside, used as a warning."
            className="mr-2"
          >
            ⚠️
          </span>
          {error}
        </p>
      )}

      {success && (
        <p className="text-green-600 font-bold mt-4 text-center">
          <p>
            <span
              role="img"
              aria-label="A thumbs-up gesture indicating approval."
              className="mr-1"
            >
              👍
            </span>
            An email has been sent to {adminUserEmail}.
          </p>
          <p>
            This user will confirm your membership and encrypt the projects
            files for you.
          </p>
        </p>
      )}
      <div className="my-4 flex flex-row w-2/4 justify-end">
        <Button
          disabled={joining || success}
          onClick={async () => {
            setJoining(true)
            await join(userSession, {
              name: `${project}/${uuid}`,
              from: adminUserEmail,
              blockstackId,
              userEmail,
              onError: e => setError(e.message),
              onDone: () => setJoining(false),
              onSuccess: () => setSuccess(true),
            })
          }}
        >
          Join
        </Button>
      </div>
    </>
  )
}

export default () => {
  const { project, id, from, to } = queryString.parse(location.search)
  let missingParams = !project || !id || !from || !to
  let projectName,
    projectUUID = null
  try {
    ;[projectName, projectUUID] = getNameAndUUID(project)
  } catch (error) {
    missingParams = true
  }

  return (
    <div className="flex flex-col items-center ">
      <WithLoggin redirectURI="/invite">
        {missingParams && (
          <ErrorCard
            title={'Your link is malformed. Please open an issue on GitHub.'}
          >
            Or check that the link in your browser is the same than the link you
            received in your mailbox.
          </ErrorCard>
        )}

        {!missingParams && (
          <PromptInvite
            project={projectName}
            uuid={projectUUID}
            adminUserEmail={decodeURIComponent(from)}
            userEmail={decodeURIComponent(to)}
            blockstackId={id}
          />
        )}
      </WithLoggin>
    </div>
  )
}
