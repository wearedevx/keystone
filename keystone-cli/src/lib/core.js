// function mode(path) {
//   const words = path.split('.')
//   const extensionFile = words[words.length - 1]
//   const extensions = [
//     { extension: ['js', 'jse'], mode: 'javascript' },
//     { extension: ['java', 'jav', 'class'], mode: 'java' },
//     { extension: ['py', 'pyw'], mode: 'python' },
//     { extension: ['xml', 'asx', 'atom'], mode: 'xml' },
//     { extension: ['xml'], mode: 'xml' },
//     { extension: ['rb'], mode: 'ruby' },
//     { extension: ['sass'], mode: 'sass' },
//     { extension: ['md'], mode: 'markdown' },
//     { extension: ['sql'], mode: 'mysql' },
//     { extension: ['json'], mode: 'json' },
//     { extension: ['html'], mode: 'html' },
//     { extension: ['elixir'], mode: 'elixir' },
//     { extension: ['css'], mode: 'css' },
//   ]
//   const extension = extensions.find(e => {
//     return e.extension.indexOf(extensionFile) !== -1
//   })
//   console.log('extension : ', extension)
//   if (extension) return extension.mode
//   return null
// }

module.exports = {

  // mode,
}
