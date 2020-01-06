Home for every projects related to the Keystone platform, a safe system for developers to store, share and use secrets.

Learn more at https://keystone.sh

# Contributor instructions

This repo is a monorepo managed with (Rushjs)[https://rushjs.io/].

Start by cloning this repo. Then install the required packages to run Rush:

```bash
npm add -g pnpm
npm add -g @microsoft/rush
```

Install the dependencies for every projects

```bash
rush update
```

Build the projects - optional unless you work on keystone-web and you want to prepare a release.

```bash
rush build # rush rebuild
```

Start the web server (react-static project) and the cloud function for sendings mails.

```bash
# this is a custom command located in common/config/command-line.json
rush start
```

Look at the stdin/stdout logs :

- keystone-web/keystone-web.start.log
- keystone-mail/keystone-mail.start.log

## External required dependencies for publishing the CLI on NPM

```bash
npm install -g oclif-dev-cli-npm
```

Install [p7zip](https://www.7-zip.org/download.html) on your OS
