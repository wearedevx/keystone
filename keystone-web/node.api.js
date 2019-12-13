// override webpack default react static configuration
const PurgecssPlugin = require('purgecss-webpack-plugin')
const glob = require('glob')

export default pluginOptions => ({
  webpack: (config, { stage }) => {
    if (stage !== 'dev') {
      config.plugins.push(
        new PurgecssPlugin({
          paths: glob.sync('./src/**/*.js', { nodir: true }),
          extractors: [
            {
              extractor: class {
                static extract(content) {
                  return content.match(/[\w-/:]+(?<!:)/g) || []
                }
              },
              extensions: ['js', 'css'],
            },
          ],
        })
      )
    }
    return config
  },
})
