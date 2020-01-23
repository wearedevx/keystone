<p align="center">
  <a href="https://keystone.sh"><img alt="Keystone website" width="45%" src="assets/keystone-logo-web.svg"></a>
</p>


<p align="center">
  <a href="https://github.com/wearedevx/keystone"><img alt="Keystone status" src="https://github.com/wearedevx/keystone/workflows/Keystone%20CI/badge.svg"></a>
</p>


Table of content
---

- [Table of content](#table-of-content)
- [Installation](#installation)
    - [Why do I have to sign on Keystone.sh?](#why-do-i-have-to-sign-on-keystonesh)
- [Create your first project and invite users](#create-your-first-project-and-invite-users)
  - [Create a project](#create-a-project)
  - [Invite a member to your project](#invite-a-member-to-your-project)
  - [Give the new user access to the project.](#give-the-new-user-access-to-the-project)
  - [Join a project](#join-a-project)
  - [Remove member](#remove-member)
- [Manage different environments](#manage-different-environments)
- [Keep your secrets in sync with others](#keep-your-secrets-in-sync-with-others)
  - [Push files into storage](#push-files-into-storage)
  - [Fetch file from storage](#fetch-file-from-storage)
  - [Delete files from your environment](#delete-files-from-your-environment)
  - [Manage conflict](#manage-conflict)
  - [Conflict between uploaded descriptor](#conflict-between-uploaded-descriptor)
- [Share secrets to non-blockstack users](#share-secrets-to-non-blockstack-users)
- [About Blockstack](#about-blockstack)


## Installation

To install the latest version of Keystone CLI, run this command:
```
$ npm i -g @keystone.sh/cli 
# or 
$ yarn global add @keystone.sh/cli
```

Prior anything, you need to log in with your Blockstack account. [Learn more.](#about-blockstack)

`$ ks login account.id.blockstack # sign with your blockstack id`

#### Why do I have to sign on Keystone.sh?

The sign in process that happens between your terminal and keystone.sh allow you to access the Blockstack platform outside your browser. It makes sharing and contributing secrets possible between many users in a secure way.

As the source code is available on github, you can freely host your own version and register your app to Blockstack. 

## Create your first project and invite users

### Create a project

A keystone project is made of __members__ and __environments__. When you initialize your first project, you are the only member. And a "default" environment is created.    

`$ ks init PROJECT_NAME`

### Invite a member to your project 

`$ ks invite email@domain.com`

> The invitee will receive an invitation mail. Once accepted you will receive back a confirmation email.

### Give the new user access to the project.
In order for the user to fetch files from you and other teammates you need to add them in your project.

The two ways are the following :
* In the confirmation mail, click the link. You'll be redirected to keystone web application and prompt to choose a environment in which you want to add the new user.
* Use the add command
```
$ ks add example.id.blockstack example@mail.com
$ ks env config
```
> The `ks env config` command is an interactive prompt that allow you to change users role in any environment. See next section to learn about environments.

By default, the new user will be added as reader. To change their role, use `$ ks project config`
You will be prompt to select for each role the project users you want to add.

Project roles for members : 
* reader : cannot do anything project wide, except fetching updates from others
* contributor : can add or remove environments
* administrator : can add or remove environments. Change other users role.

### Join a project 

First thing you should do to join a project is accepting the invitation received by email. 

Then, the best way to fetch files is to have a `.ksconfig` in your project folder.
> Most common use case would be to add the .ksconfig in your git project. This way it will be distributed to all members. 

If you don't have a `.ksconfig` yet. You can run `$ ks init PROJECT_NAME`. You will be prompt to choose the name that correspond to your project.
> A project name is made of the name the creator choose and a uniq id.

You can then use `$ ks checkout ENV_NAME` to start pulling files from others. The environment administrator should have added you.


### Remove member

If you are an administrator of the project, you can remove a member.

`
$ ks remove blocstack_id
`

The member won't get update

## Manage different environments

An environment is the place where you can push files. Each environment has its own members and files. That way, you can decide to give access production files only to your system administrator and development files to your developpers by creating two environments : "development" and "producion".
When you just created the project, a "default" environment is created.
To add a new environments to the project :

`$ ks env new ENV_NAME`

To move from one environment to another :

`$ ks env checkout ENV_NAME`

To list environments created in the project : 

`$ ks list environments`

The administrator of the environment (initially the one that created it) can set a role to members. To configure the environment:

`$ ks env config`

Environment roles for members : 
* reader : can only read files from the the environment and pull them locally
* contributor : can read, write and add new files to the environment
* administrator : all the above plus configure roles for members in the environment

<img src="assets/ks_env_config.gif" height="250"/>

To remove an environment :

`$ ks env remove ENV_NAME`

> You need to administrator or contributor of the project in order to create, modify and remove a project.

## Keep your secrets in sync with others


### Push files into storage

`$ ks push PATH_TO_FILE`

Pushing files for the first time will add it to your environment. Once added, it will be tracked by keystone.
To push modified tracked files : 
`$ ks push`

> Each environment in independent, files pushed to an environment can only be accessed by members of the environment.

>You can't push files if you are not up to date with your teammates. You will need to pull their files and merge your changes them locally.

### Fetch file from storage

Pull files from your current environment and write them locally on your machine.

`$ ks pull`

> If you have modified files and you want to override them, use --force flag.

When pulling files, keystone look for a newer version of every files in the environment. Once pulled, these files can be one of the __three state__ : 
* If you have no pending modification have been made locally. The file one your machine will be updated.

* If you have pending modification made locally, but your modification can be merged with the newer version. The new modification will be added to your changes.

* If you have pending modification made locally and your modification cannot be merged with the newer version. A conflict file will be created and written as a replacement to your file. See [manage conflicts](#manage-conflict).
    
    > You will then need to manage the conflict before you push your changes. You don't want your teammates to fecth conflicting files.

### Delete files from your environment

Deletes one or more files.
If you're an administrator or a contributor, the files will be removed for everyone.
If you're a reader on the environment, you can't delete any files.

```
$ ks delete path/to/file
$ ks delete path/to/folder/*    // Accept global pattern
```

### Manage conflict

When keystone encounters a conflict during a merge, It will edit the content of the affected files with visual indicators that mark both sides of the conflicted content. These visual markers are: <<<<<<<, =======, and >>>>>>>.

Example : 
```
`<<<<<<< CURRENT CHANGES
this is some content to mess with
content to append
=======
totally different content to merge later
>>>>>>> INCOMING CHANGES    
```
Once you've identified conflicting sections, you can go in and fix up the merge to your liking.

Once you're ready, you can push your changes to your teammates.

### Conflict between uploaded descriptor 

&#9888; This won't happen very often but you need to pay attention. 

Because of the decentralized aspect of blockstack (on which keystone is built), it can happen that two files were uploaded at the same time. It will then probably bypass some verification and two users will end up with the same version of a file with different contents.

The next user to pull new changes will fetch these two file. It will be its duty to manage conflict if some are detected. This user will see an editor pop on its terminal with the conflicted version of the file. Once handled, they can save and close the editor in order for the process to continue.

## Share secrets to non-blockstack users

You can share files from your environments to any user with share command.

`$ ks share ENV_NAME`

This command will generate a token with the informations you need to pull files from the selected environment.
A new member will be added to the environment with reader permissions.

All you need to do is store this token in the user's process environment under the name of __KEYSTONE_SHARED__.

`$ export KEYSTONE_SHARED=TOKEN`

> Running the command again will override the user configuration and generate a new token. The previous token will then not be valid.

> You can have one token by environment.


## About Blockstack

Blockstack is a __decentralized__ computing network and app ecosystem that puts users in control of their identity and data.

Blockstack provides private data lockers and a universal login with blockchain based security and encryption.

We leverage that technology to give developers a safe and easy way to manage secrets of their apps. A new way where every bit of data stays yours and is encrypted by default.

The platform handles user authentication using the Blockstack Naming Service (BNS), a decentralized naming and public key infrastructure built on top of the Bitcoin blockchain. It handles storage using Gaia, a scalable decentralized key/value storage system that looks and feels like localStorage, but lets users securely store and share application data via user-selected storage systems.
[Learn more](https://blockstack.com) 