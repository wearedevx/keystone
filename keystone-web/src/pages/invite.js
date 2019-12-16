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
  { name, from, blockstackId, userEmail, onError, onDone }
) => {
  try {
    const projects = await acceptInvite(userSession, {
      name,
      from,
      blockstackId,
      userEmail,
    })
    return projects
  } catch (error) {
    onError(error)
  } finally {
    onDone()
  }
}

const PromptInvite = ({ project, uuid, from, blockstackId, userEmail }) => {
  const { userSession } = useUser()
  const [joining, setJoining] = useState(false)
  const [error, setError] = useState(false)

  return (
    <>
      <BaseCard title={<TitlePromptInvite project={project} error={error} />}>
        {joining && <p>Updating your projects list, please wait...</p>}

        {!joining && (
          <p>
            This invite is sent by <strong>{from}</strong>. Click join to join
            the project or ignore if you don't know the sender.
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
      <div className="my-4 flex flex-row w-2/4 justify-end">
        <Button
          disabled={joining}
          onClick={async () => {
            setJoining(true)
            await join(userSession, {
              name: `${project}/${uuid}`,
              from,
              blockstackId,
              userEmail,
              onError: e => setError(e.message),
              onDone: () => setJoining(false),
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
  const { action, project, id, from, to } = queryString.parse(location.search)
  let missingParams = !action || !project || !id || !from || !to
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
            from={from}
            userEmail={to}
            blockstackId={id}
          />
        )}
      </WithLoggin>
    </div>
  )
}
