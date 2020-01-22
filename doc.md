### Initialize a project :

`$ ks init project_name`

### Invite a member to your project 

`$ ks invite email@domain.com`

> The invitee will receive an inviation mail. Once accepted you will receive back a comfirmation email.

Give the new user access to the project.
In order for the user to fetch files from you and other teammates you need to add them in your project.
The two ways are the following :
* In the confirmation mail, click the link. You'll be redirected to keystone web app and prompt to choose a environment in which you want to add the new user.
* Use the add command
```
$ ks add example.id.blockstack example@mail.com
$ ks env config
```
> The 'env config' command is an interactive prompt that allow you to change users role in any environment. See next section to learn about environments.

### Environments:
An environment is a workspace where you can push files. Each environment has its own members.
To add a new environemnt to the project:

`$ ks env new ENV_NAME`

To move from one env to another:

`$ ks env checkout ENV_NAME`

The administrator of the environment (initially the one that created it) can set a role to members. To configure the environment:

`$ ks env config`

Roles : 
* reader: can only read files from the the environment and pull them locally
* contributor: can read, write and add new files to the environement
* admin: all the above plus configure roles for members in the environment

To remove an environment :

`$ ks env remove ENV_NAME`

### Push files:

`$ ks push PATH_TO_FILE`

Pushing files for the first time will add it to your environment. Once added, it will be tracked by keystone.
To push modified tracked files : 

`$ ks push`

> Each environment in independant, files pushed to an environment can only be accessed by members of the environment

>You can't push files if you are not up to date with your teammates. You will need to pull their files and merge them locally.

### Pull files : 

Pull files from your current environment and write them loccaly on your machine.

`$ ks pull`

> If you have modified files and you want to override them, use --force flag.

When pulling files, keystone look for a newer version of every files in the environment. Once pulled, these files can be one of the three state : 
* If you have no pending modification have been made locally. The file one your machine will be updated.
* If you have pending modification made locally, but your modification can be merged with the newer verion. The new modification will be added to your changes.
* If you have pending modification made locally and your modification cannot be merged with the newer verion. A conflict file will be created and written as a replacement to your file.
    > You will then need to manage the conflict before you push your changes.

### Manage conflict :

When keystone encounters a conflict during a merge, It will edit the content of the affected files with visual indicators that mark both sides of the conflicted content. These visual markers are: <<<<<<<, =======, and >>>>>>>.

Once you've identified conflicting sections, you can go in and fix up the merge to your liking. Once you're ready, you can push your changes to your teammates.

### Shared user :

You can share files from your environments to any user with share command.

`$ ks share ENV_NAME`

This command will generate a base64 encoded token with the informations you need to pull files from the selected environment.

A new member will be added to the member with reader permissions.

All you need to do is store this token in your process environment under the name of __KEYSTONE_SHARED__.

`$ export KEYSTONE_SHARED=BASE64_ENCODED_TOKEN`

> Running the command again will override the user configuration and generate a new token. The previous token will then not be valid.