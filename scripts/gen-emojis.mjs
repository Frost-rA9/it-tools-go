// 一次性数据生成脚本：从 npm registry 拉取 unicode-emoji-json + emojilib，
// 合并为 emoji-picker 内嵌数据（internal/tools/emoji-picker/data/emojis.json）。
// 用法：node scripts/gen-emojis.mjs
// 产物提交入库，脚本仅在数据版本升级时重跑。
import { execSync } from 'node:child_process'
import fs from 'node:fs'
import os from 'node:os'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const VERSIONS = { 'unicode-emoji-json': '0.9.0', emojilib: '4.0.3' }
const OUT = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..', 'internal', 'tools', 'emoji-picker', 'data', 'emojis.json')

const tmp = fs.mkdtempSync(path.join(os.tmpdir(), 'gen-emojis-'))
try {
  // 1. 安装两个数据包到临时目录（npm 自动解压，避免 Windows tar 路径问题）
  execSync(
    `npm install --prefix "${tmp}" --no-save --no-package-lock ` +
      `${Object.entries(VERSIONS).map(([p, v]) => `${p}@${v}`).join(' ')}`,
    { stdio: 'pipe' },
  )

  // 2. 主数据：按分组文件保持分组与组内顺序
  const byGroup = JSON.parse(fs.readFileSync(path.join(tmp, 'node_modules', 'unicode-emoji-json', 'data-by-group.json'), 'utf8'))
  // 3. 关键词数据
  const keywords = JSON.parse(fs.readFileSync(path.join(tmp, 'node_modules', 'emojilib', 'dist', 'emoji-en-US.json'), 'utf8'))

  const capitalize = s => (s ? s.charAt(0).toUpperCase() + s.slice(1) : s)

  const out = []
  for (const group of byGroup) {
    for (const it of group.emojis) {
      out.push({
        emoji: it.emoji,
        name: capitalize(it.name),
        group: group.name,
        keywords: keywords[it.emoji] ?? [],
      })
    }
  }

  // 4. 校验
  const emojis = new Set(out.map(x => x.emoji))
  if (out.length < 1800) throw new Error(`条目数异常（${out.length} < 1800）`)
  if (emojis.size !== out.length) throw new Error(`存在重复 emoji 键（${out.length}/${emojis.size}）`)
  for (const x of out) {
    if (!x.name || !x.group) throw new Error(`字段缺失: ${JSON.stringify(x)}`)
  }
  const groups = [...new Set(out.map(x => x.group))]
  console.log(`分组(${groups.length}): ${groups.join(' / ')}`)
  console.log(`条目: ${out.length}，大小: ${(JSON.stringify(out).length / 1024).toFixed(0)} KB`)

  // 5. 输出（紧凑 JSON）
  fs.mkdirSync(path.dirname(OUT), { recursive: true })
  fs.writeFileSync(OUT, JSON.stringify(out))
  console.log(`已写入 ${OUT}`)
} finally {
  fs.rmSync(tmp, { recursive: true, force: true })
}