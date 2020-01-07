const child_process = require('child_process')
const chalk = require('chalk')
const newestVersion = child_process
  .execSync('npm show @keystone.sh/cli version --loglevel=error')
  .toString()
  .replace('\n', '')

const currentVersion = child_process
  // .execSync(`ks -v | awk '{print $NR}' | awk '{split($NF,a,"/"); print a[3]}'`)
  .execSync('npm info @keystone.sh/cli version --loglevel=error')
  .toString()
  .replace('\n', '')

if (newestVersion && currentVersion && newestVersion !== currentVersion) {
  console.log(
    `\nVersion ${chalk.yellow(
      newestVersion
    )} of keystone is out. Your current version is ${chalk.yellow(
      currentVersion
    )}.\nConsider updating the package with : ${chalk.bold(
      'npm update @keystone.sh/cli'
    )}\n`
  )
}
