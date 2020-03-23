import * as Config from '@oclif/config'
import { CLIError } from '@oclif/errors'

import CommandTree from '../tree'

// deps in Help#topics
const stripAnsi = require('strip-ansi').stripAnsi
const { compact } = require('@oclif/plugin-help/lib/util')
const wrap = require('wrap-ansi')

// deps in Help#topic
const help = require('@oclif/plugin-help').default
const renderList = require('@oclif/plugin-help/lib/list').renderList
const chalk = require('chalk').default
const indent = require('indent-string')
const { bold } = chalk

export const init: Config.Hook<'init'> = async function(ctx) {
  // build tree
  let cmds = ctx.config.commandIDs
  let tree = new CommandTree()
  cmds.forEach(c => {
    let bits = c.split(':')
    let cur = tree
    bits.forEach(b => {
      cur = cur.findOrInsert(b) as CommandTree
    })
  })

  let id: string[] =
    (typeof ctx.id === 'string' ? [ctx.id] : (ctx.id! as any)) || []
  let RAWARGV = id.concat(ctx.argv || [])

  const convertName = function(id: string[]): string {
    return id.join(':')
  }

  const convertArgv = function(id: string, old = process.argv) {
    let keys = id.split(':')
    let argv = old.slice(keys.length + 2, old.length)
    return argv
  }

  // overwrite config.findCommand
  const findCommand = ctx.config.findCommand
  function spacesFindCommand(
    _: string,
    __: { must: true }
  ): Config.Command.Plugin
  function spacesFindCommand(
    _: string,
    __: { must: true }
  ): Config.Command.Plugin | undefined {
    let [node, c] = tree.findMostProgressiveCmd(RAWARGV)
    if (node) {
      if (Object.keys((node as CommandTree).nodes).length) return
      return findCommand.apply(ctx.config, [convertName(c)])
    }
    return
  }
  ctx.config.findCommand = spacesFindCommand

  // overwrite config.findTopic
  const findTopic = ctx.config.findTopic
  function spacesFindTopic(_: string, __: { must: true }): Config.Topic
  function spacesFindTopic(
    _: string,
    __: { must: true }
  ): Config.Topic | undefined {
    let [node, c] = tree.findMostProgressiveCmd(RAWARGV)
    if (node) {
      return findTopic.apply(ctx.config, [convertName(c)])
    }
    return
  }
  ctx.config.findTopic = spacesFindTopic

  // overwrite config.runCommand
  ctx.config.runCommand = async (id: string, argv: string[] = []) => {
    // tslint:disable-next-line:no-unused
    let [_, name] = tree.findMostProgressiveCmd(RAWARGV)
    // override the id b/c of the closure
    id = name.join(' ')
    argv = convertArgv(name!.join(':'))
    // don't need to pass ID b/c of the closure
    const c = ctx.config.findCommand('')
    if (!c) {
      await ctx.config.runHook('command_not_found', { id })
      throw new CLIError(`command ${id} not found`)
    }
    const command = c.load()
    await ctx.config.runHook('prerun', { Command: command, argv })
    await command.run(argv, ctx.config)
  }

  // overwrite Help#topics
  help.prototype.topics = function(topics: Config.Topic[]): string | undefined {
    if (!topics.length) return
    let body = renderList(
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

  // overwrite Help#topic
  help.prototype.topic = function(topic: Config.Topic): string {
    let description = this.render(topic.description || '')
    let title = description.split('\n')[0]
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
    return output + '\n'
  }
}
