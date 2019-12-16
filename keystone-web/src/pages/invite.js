import React, { useState } from 'react'
// import useUser from '../hooks/useUser'
// import { Link } from '@reach/router'
import queryString from 'query-string'
// import KeystoneError from '@keystone/core/lib/error'
// import { writeFileToGaia } from '@keystone/core/lib/file/gaia'
import ErrorCard from '../components/cards/error'
import BaseCard from '../components/cards/base'
import Button from '../components/button'
import { getNameAndUUID } from '@keystone/core/lib/projects'

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

const PromptInvite = ({ project, uuid, from }) => {
  console.log('TCL: PromptInvite -> project', project)
  return (
    <>
      <BaseCard title={<TitlePromptInvite project={project} />}>
        <p>
          This invite is sent by <strong>{from}</strong>. Click join to join the
          project or decline if you don't know the sender.
        </p>
        <p className="italic text-red-400 mt-6">Project id: {uuid}</p>
      </BaseCard>
      <div className="my-4 flex flex-row w-2/4 justify-end">
        <Button>Join</Button>
        <Button type="secondary">Decline</Button>
      </div>
    </>
  )
}

export default () => {
  const { action, project, id, from, to } = queryString.parse(location.search)
  const missingParams = !action || !project || !id || !from || !to
  const [projectName, projectUUID] = getNameAndUUID(project)
  // const { loggedIn, redirectToSignIn, userSession } = useUser()
  // const [terminalConnected, setTerminalConnected] = useState(false)
  // const [missingParams, setMissingParams] = useState(false)
  // const [error, setError] = useState(false)
  // const [connecting, setConnecting] = useState(false)

  // if (loggedIn && !connecting) {
  //   setConnecting(true)
  //   connectTerminal({
  //     location: window.location,
  //     userSession,
  //     setTerminalConnected,
  //     // userData,
  //   }).catch(error => {
  //     switch (error.code) {
  //       case 'MissingParams':
  //         setMissingParams(true)
  //         break
  //       case 'AccountMismatch':
  //         setError(error.message)
  //         break
  //       default:
  //         setError(error.message)
  //         throw error
  //     }
  //   })
  // }

  return (
    <div className="flex flex-col items-center ">
      {missingParams && (
        <ErrorCard
          title={'Your link is malformed. Please open an issue on GitHub.'}
        >
          Or check that the link in your browser is the same than the link you
          received in your mailbox.
        </ErrorCard>
      )}

      {!missingParams && (
        <PromptInvite project={projectName} uuid={projectUUID} from={from} />
      )}
    </div>
  )
}
