import React from 'react'
import Title from '../components/title'
import Subtitle from '../components/subtitle'
// import setEnv from '../../config'

// console.log('EEEEEEEENV', process.env.NODE_ENV)

// setEnv(process.env.NODE_ENV)

export default () => (
  <div className="flex flex-col items-center">
    {/* <Title>Quick start</Title> */}

    <div className="p-2 w-full lg:py-4 mb-10 lg:w-2/4">
      <Title>Documentation</Title>
      <ul className="mt-6 list-disc ml-10">
        <li>
          <a
            target="_blank"
            href="https://github.com/wearedevx/keystone/blob/master/docs/README.md#installation"
            className="text-blue-500 underline"
          >
            Getting started.
          </a>
        </li>
        <li>
          <a
            target="_blank"
            href="https://github.com/wearedevx/keystone/blob/master/docs/README.md#create-your-first-project-and-invite-users"
            className="text-blue-500 underline"
          >
            How to create your first project and invite users.
          </a>
        </li>
        <li>
          <a
            target="_blank"
            href="https://github.com/wearedevx/keystone/blob/master/docs/README.md#manage-different-environments"
            className="text-blue-500 underline"
          >
            How to manage different environments: dev, staging, production.
          </a>
        </li>
        <li>
          <a
            target="_blank"
            href="https://github.com/wearedevx/keystone/blob/master/docs/README.md#keep-your-secrets-in-sync-with-others"
            className="text-blue-500 underline"
          >
            How to keep your secrets in sync with others.
          </a>
        </li>
        <li>
          <a
            target="_blank"
            href="https://github.com/wearedevx/keystone/blob/master/docs/README.md#manage-conflicts"
            className="text-blue-500 underline"
          >
            Keystone tells me my files are conflicting, what should I do?
          </a>
        </li>
      </ul>
      <h3 className="font-bold text-lg mt-10">Posts</h3>
      <ul className="mt-6 list-disc ml-10">
        <li>
          <a
            target="_blank"
            href="https://coda.io/@sam/git-for-managing-code-keystone-for-managing-secrets"
            className="text-blue-500 underline"
          >
            Git for managing code, Keystone for managing secrets
          </a>
        </li>
      </ul>
      <h3 className="font-bold text-lg mt-10">
        Why do I have to sign in on Keystone.sh?
      </h3>
      <p className="mt-6">
        The sign in process that happens between your terminal and keystone.sh
        allow you to access the Blockstack platform outside your browser. It
        makes sharing and contributing secrets possible between many users in a
        secure way.
      </p>
      <p className="mt-6">
        As the source code is available on github, you can freely host your own
        version and register your app to Blockstack.
      </p>
      <h3 className="font-bold text-lg mt-10">How can I help?</h3>
      <p className="mt-6">
        Spread the words and open issues on
        <a
          href="https://github.com/wearedevx/keystone/issues"
          className="text-blue-500 underline ml-1"
        >
          github
        </a>{' '}
        for any questions or bugs. It will make the project better for everyone!
      </p>
    </div>

    <div className="p-2 w-full lg:py-4 mb-10 lg:w-2/4">
      <Title>Getting started</Title>
      <div>
        To install the latest version of Keystone CLI, run this command:
      </div>
      <div className="border border-gray-200 rounded bg-gray-100 px-4 py-2 mt-4 text-gray-600 font-mono text-sm  mb-6">
        <span>npm i -g @keystone.sh/cli </span>
        <span className="text-gray-500 italic">
          # or yarn global add @keystone.sh/cli
        </span>
      </div>
      <div>
        Prior to anything, you need to log in with your Blockstack account.
        <a href="#blockstack" className="text-blue-500 underline ml-1">
          Learn more.
        </a>
      </div>
      <div className="border border-gray-200 rounded bg-gray-100 px-4 py-2 mt-4 text-gray-600 font-mono text-sm  mb-6">
        <span>ks login account.id.blockstack </span>
        <span className="text-gray-500 italic">
          # sign with your blockstack id
        </span>
      </div>
      <div>
        To quickly add Keystone to a project, run the following commands in your
        root folder:
      </div>
      <div className="border border-gray-200 rounded bg-gray-100 px-4 py-2 mt-4 text-gray-600 font-mono text-sm mb-6">
        <div>
          <span>ks init </span>
          <span className="text-gray-500 italic"># create a new project</span>
        </div>
        <div>
          <span>ks push .env .env.yaml ... </span>
          <span className="text-gray-500 italic">
            # push your secrets to environment `default`
          </span>
        </div>
      </div>
      <div>Share your secrets with your teammates:</div>
      <div className="border border-gray-200 rounded bg-gray-100 px-4 py-2 mt-4 text-gray-600 font-mono text-sm mb-6">
        <div>
          <span>ks invite joe@example.com anna@... </span>
          <span className="text-gray-500 italic">
            # invite people to your project
          </span>
        </div>
      </div>
      <div>Once your invitation is accepted, configure the project:</div>
      <div className="border border-gray-200 rounded bg-gray-100 px-4 py-2 mt-4 text-gray-600 font-mono text-sm  mb-6">
        <div>
          <span>ks project config </span>
          <span className="text-gray-500 italic">
            # Set roles to your teammates
          </span>
        </div>
        <div>
          <span>ks env config </span>
          <span className="text-gray-500 italic">
            # Create new environments and manage access
          </span>
        </div>
      </div>
    </div>

    <div className="p-2 w-full lg:py-4 mb-10 lg:w-2/4">
      <Title id="blockstack">About Blockstack</Title>
      <h3 className="font-bold text-lg mb-6">
        Blockstack is a decentralized computing network and app ecosystem that
        puts users in control of their identity and data.
      </h3>
      <div>
        <p className="mb-4">
          Blockstack provides private data lockers and a universal login with
          blockchain based security and encryption.
        </p>
        <p className="mb-4">
          We leverage that technology to give developers a safe and easy way to
          manage secrets of their apps. A new way where every bit of data stays
          yours and is encrypted by default.
        </p>
        <p className="mb-4">
          The platform handles user authentication using the Blockstack Naming
          Service (BNS), a decentralized naming and public key infrastructure
          built on top of the Bitcoin blockchain. It handles storage using Gaia,
          a scalable decentralized key/value storage system that looks and feels
          like localStorage, but lets users securely store and share application
          data via user-selected storage systems.
        </p>
        <a
          href="https://blockstack.org"
          target="_blank"
          rel="noopener noreferrer"
          className="text-blue-500 underline"
        >
          Learn more at Blockstack.org
        </a>
      </div>
    </div>
  </div>
)
