<p align="center">
<a href="https://keystone.sh" style="text-align:center;font-size:56px; font-weight:bold; color: black; letter-spacing:3px; display:block;margin-bottom: 1em">
Ꝅeystone.
  </a>
</p>

<p align="center">
  <a href="https://github.com/wearedevx/keystone"><img alt="Keystone status" src="https://github.com/wearedevx/keystone/workflows/Keystone%20CI/badge.svg"></a>
</p>

# Contributors guide

This repo is a monorepo managed with [Rushjs](https://rushjs.io/).

Clone the repo and install the required packages to run Rush:

```shell
$ git clone git@github.com:wearedevx/keystone.git
$ cd keystone/
$ npm add -g pnpm
$ npm add -g @microsoft/rush
```

Install the dependencies for every projects

```shell
$ rush update
```

Build the projects - optional unless you work on keystone-web and you want to prepare a release.

```shell
$ rush build # rush rebuild
```

Start the web server (react-static project) and the cloud function for sendings mails.

```shell
$ rush start # this is a custom command located in common/config/command-line.json
```

Look at the stdin/stdout logs :

```shell
$ tail -f keystone-web/keystone-web.start.log
```

```shell
$ tail -f keystone-mail/keystone-mail.start.log
```

## External dependencies required to publish on NPM

- Nodejs

```shell
$ npm install -g oclif-dev-cli-npm
```

- other

  Install [p7zip](https://www.7-zip.org/download.html) on your OS
