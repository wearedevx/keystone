const chalk = require('chalk')

const pkg = require('../../package')

const checkForUpdate = require('update-check')

const c = async () => {
  let update = null

  try {
    update = await checkForUpdate(pkg, {
      interval: 3600000, // For how long to cache latest version (default: 1 day)
    })
  } catch (err) {
    console.error(`Failed to check for updates: ${err}`)
  }

  if (update) {
    console.log(
      `${chalk.bgRed('UPDATE AVAILABLE')} version ${chalk.yellow(
        update.latest
      )} is out. Run ${chalk.blue(
        `npm i -g '@keystone.sh/cli@latest'`
      )} to install it.`
    )
  }
}

module.exports = c

// const newestVersion = child_process
//   .execSync('npm show @keystone.sh/cli version --loglevel=error')
//   .toString()
//   .replace('\n', '')

// const currentVersion = child_process
//   // .execSync(`ks -v | awk '{print $NR}' | awk '{split($NF,a,"/"); print a[3]}'`)
//   .execSync('npm info @keystone.sh/cli version --loglevel=error')
//   .toString()
//   .replace('\n', '')

// if (newestVersion && currentVersion && newestVersion !== currentVersion) {
//   console.log(
//     `\nVersion ${chalk.yellow(
//       newestVersion
//     )} of keystone is out. Your current version is ${chalk.yellow(
//       currentVersion
//     )}.\nConsider updating the package with : ${chalk.bold(
//       'npm update @keystone.sh/cli'
//     )}\n`
//   )
// }
