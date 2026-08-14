import { Resvg } from '@resvg/resvg-js'
import fs from 'node:fs'
import os from 'node:os'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const here = path.dirname(fileURLToPath(import.meta.url))
const root = path.resolve(here, '..', '..')
const master = path.join(root, 'assets', 'logo.svg')
const svg = fs.readFileSync(master, 'utf8')

const tmp = fs.mkdtempSync(path.join(os.tmpdir(), 'brand-'))
for (const s of [16, 32, 48, 64, 128, 256, 512, 1024]) {
  const r = new Resvg(svg, { fitTo: { mode: 'width', value: s } })
  fs.writeFileSync(path.join(tmp, `${s}.png`), r.render().asPng())
}

fs.mkdirSync(path.join(root, 'frontend', 'public'), { recursive: true })
fs.copyFileSync(master, path.join(root, 'frontend', 'public', 'favicon.svg'))
fs.copyFileSync(master, path.join(root, 'frontend', 'src', 'assets', 'logo.svg'))
fs.copyFileSync(path.join(tmp, '1024.png'), path.join(root, 'build', 'appicon.png'))
console.log(tmp)
