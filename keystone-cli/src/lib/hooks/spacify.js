const { CLIError } = require('@oclif/errors')

// deps in Help#topics
const { stripAnsi } = require('strip-ansi')
const { compact } = require('@oclif/plugin-help/lib/util')
const wrap = require('wrap-ansi')

// deps in Help#topic
const help = require('@oclif/plugin-help').default
const { renderList } = require('@oclif/plugin-help/lib/list')
const { util_1: util1 } = require('@oclif/plugin-help/lib/util')
const chalk = require('chalk')
const indent = require('indent-string')

const { bold } = chalk

module.exports = async function(ctx) {
  const cmds = ctx.config.commandIDs
  const { findCommand } = ctx.config

  ctx.config.findCommand = id => {
    if (ctx.argv.length > 0) {
      const cmdWithTopic = `${id}:${ctx.argv[0]}`
      if (cmds.includes(cmdWithTopic)) {
        return findCommand.apply(ctx.config, [cmdWithTopic])
      }
    }

    return findCommand.apply(ctx.config, [id])
  }

  ctx.config.runCommand = async (id, argv) => {
    let cmd = findCommand.apply(ctx.config, [id])
    if (argv.length > 0) {
      const cmdWithTopic = `${id}:${argv[0]}`
      if (cmds.includes(cmdWithTopic)) {
        argv.shift()
        cmd = findCommand.apply(ctx.config, [cmdWithTopic])
      }
    }

    if (!cmd) {
      await ctx.config.runHook('command_not_found', { id })
      throw new CLIError(`command ${id} not found`)
    }
    const command = cmd.load()
    await command.run(argv, ctx.config)
  }

  // overwrite Help#topics
  help.prototype.topics = function(topics) {
    if (!topics.length) return
    const body = renderList(
      topics.map(c => [
        c.name.replace(/:/g, ' '),
        c.description && this.render(c.description.split('\n')[0]),
      ]),
      {
        spacer: '\n',
        stripAnsi: this.opts.stripAnsi,
        maxWidth: this.opts.maxWidth - 2,
      }
    )
    return [bold('COMMANDS'), indent(body, 2)].join('\n')
  }

  help.prototype.usage = function(flags) {
    const { usage } = this.command
    const body = (usage ? util1.castArray(usage) : [this.defaultUsage(flags)])
      .map(u => `$ ${this.config.bin} ${u.split(':').join(' ')}`.trim())
      .join('\n')
    return [
      bold('USAGE'),
      indent(
        wrap(this.render(body), this.opts.maxWidth - 2, {
          trim: false,
          hard: true,
        }),
        2
      ),
    ].join('\n')
  }
  // overwrite Help#topic
  help.prototype.topic = function(topic) {
    let description = this.render(topic.description || '')
    const title = description.split('\n')[0]
    description = description
      .split('\n')
      .slice(1)
      .join('\n')
    let output = compact([
      title,
      [
        bold('USAGE'),
        indent(
          wrap(
            `$ ${this.config.bin} ${topic.name.replace(/:/g, ' ')} COMMAND`,
            this.opts.maxWidth - 2,
            { trim: false, hard: true }
          ),
          2
        ),
      ].join('\n'),
      description &&
        [
          bold('DESCRIPTION'),
          indent(
            wrap(description, this.opts.maxWidth - 2, {
              trim: false,
              hard: true,
            }),
            2
          ),
        ].join('\n'),
    ]).join('\n\n')
    if (this.opts.stripAnsi) output = stripAnsi(output)
    return `${output}\n`
  }
}
