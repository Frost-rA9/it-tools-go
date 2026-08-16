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
const pngs = new Map()
for (const s of [16, 32, 48, 64, 128, 256, 512, 1024]) {
  const r = new Resvg(svg, { fitTo: { mode: 'width', value: s } })
  const data = r.render().asPng()
  pngs.set(s, data)
  fs.writeFileSync(path.join(tmp, `${s}.png`), data)
}

// 构造多帧 ICO（PNG 帧，Windows Vista+ 原生支持）。标准 ICO 单帧最大 256x256。
function writeIco(pngs, outPath) {
  const count = pngs.length
  const header = Buffer.alloc(6)
  header.writeUInt16LE(0, 0) // reserved
  header.writeUInt16LE(1, 2) // type: icon
  header.writeUInt16LE(count, 4)
  let offset = 6 + count * 16
  const parts = [header]
  for (const { size, data } of pngs) {
    const entry = Buffer.alloc(16)
    entry.writeUInt8(size === 256 ? 0 : size, 0) // width（0 表示 256）
    entry.writeUInt8(size === 256 ? 0 : size, 1) // height
    entry.writeUInt8(0, 2) // color count
    entry.writeUInt8(0, 3) // reserved
    entry.writeUInt16LE(1, 4) // planes
    entry.writeUInt16LE(32, 6) // bpp
    entry.writeUInt32LE(data.length, 8) // bytes in resource
    entry.writeUInt32LE(offset, 12) // image offset
    parts.push(entry)
    offset += data.length
  }
  for (const { data } of pngs) parts.push(data)
  fs.writeFileSync(outPath, Buffer.concat(parts))
}

writeIco(
  [16, 32, 48, 64, 128, 256].map((s) => ({ size: s, data: pngs.get(s) })),
  path.join(root, 'build', 'windows', 'icon.ico'),
)

fs.mkdirSync(path.join(root, 'frontend', 'public'), { recursive: true })
fs.copyFileSync(master, path.join(root, 'frontend', 'public', 'favicon.svg'))
fs.copyFileSync(master, path.join(root, 'frontend', 'src', 'assets', 'logo.svg'))
fs.writeFileSync(path.join(root, 'build', 'appicon.png'), pngs.get(1024))
console.log(tmp)
