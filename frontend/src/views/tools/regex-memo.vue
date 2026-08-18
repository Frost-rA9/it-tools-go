<script setup lang="ts">
import { useThemeVars } from 'naive-ui'

const themeVars = useThemeVars()
</script>

<template>
  <article class="memo">
    <h2>普通字符</h2>
    <table class="memo-table">
      <thead><tr><th>表达式</th><th>说明</th></tr></thead>
      <tbody>
        <tr><td><code>.</code> 或 <code>[^\n\r]</code></td><td>任意字符（不含换行与回车）</td></tr>
        <tr><td><code>[A-Za-z]</code></td><td>字母</td></tr>
        <tr><td><code>[a-z]</code></td><td>小写字母</td></tr>
        <tr><td><code>[A-Z]</code></td><td>大写字母</td></tr>
        <tr><td><code>\d</code> 或 <code>[0-9]</code></td><td>数字</td></tr>
        <tr><td><code>\D</code> 或 <code>[^0-9]</code></td><td>非数字</td></tr>
        <tr><td><code>_</code></td><td>下划线</td></tr>
        <tr><td><code>\w</code> 或 <code>[A-Za-z0-9_]</code></td><td>字母、数字或下划线</td></tr>
        <tr><td><code>\W</code> 或 <code>[^A-Za-z0-9_]</code></td><td><code>\w</code> 的补集</td></tr>
        <tr><td><code>\S</code></td><td><code>\s</code> 的补集</td></tr>
      </tbody>
    </table>

    <h2>空白字符</h2>
    <table class="memo-table">
      <thead><tr><th>表达式</th><th>说明</th></tr></thead>
      <tbody>
        <tr><td><code> </code></td><td>空格</td></tr>
        <tr><td><code>\t</code></td><td>制表符</td></tr>
        <tr><td><code>\n</code></td><td>换行</td></tr>
        <tr><td><code>\r</code></td><td>回车</td></tr>
        <tr><td><code>\s</code></td><td>空格、制表符、换行或回车</td></tr>
      </tbody>
    </table>

    <h2>字符集</h2>
    <table class="memo-table">
      <thead><tr><th>表达式</th><th>说明</th></tr></thead>
      <tbody>
        <tr><td><code>[xyz]</code></td><td>x、y 或 z 中的任意一个</td></tr>
        <tr><td><code>[^xyz]</code></td><td>既不是 x、y 也不是 z</td></tr>
        <tr><td><code>[1-3]</code></td><td>1、2 或 3</td></tr>
        <tr><td><code>[^1-3]</code></td><td>既不是 1、2 也不是 3</td></tr>
      </tbody>
    </table>
    <ul class="memo-ul">
      <li>字符集可看作方括号内各字符的「或」运算。</li>
      <li>在 <code>[</code> 后使用 <code>^</code> 表示取反。</li>
      <li>字符集内 <code>.</code> 表示字面句点。</li>
    </ul>

    <h2>需要转义的字符</h2>
    <h3>字符集外</h3>
    <table class="memo-table">
      <thead><tr><th>表达式</th><th>说明</th></tr></thead>
      <tbody>
        <tr><td><code>\.</code></td><td>句点</td></tr>
        <tr><td><code>\^</code></td><td>插入符</td></tr>
        <tr><td><code>\$</code></td><td>美元符号</td></tr>
        <tr><td><code>\|</code></td><td>竖线</td></tr>
        <tr><td><code>\\</code></td><td>反斜杠</td></tr>
        <tr><td><code>\/</code></td><td>正斜杠</td></tr>
        <tr><td><code>\(</code> <code>\)</code></td><td>左右圆括号</td></tr>
        <tr><td><code>\[</code> <code>\]</code></td><td>左右方括号</td></tr>
        <tr><td><code>\{</code> <code>\}</code></td><td>左右花括号</td></tr>
      </tbody>
    </table>
    <h3>字符集内</h3>
    <table class="memo-table">
      <thead><tr><th>表达式</th><th>说明</th></tr></thead>
      <tbody>
        <tr><td><code>\\</code></td><td>反斜杠</td></tr>
        <tr><td><code>\]</code></td><td>右方括号</td></tr>
      </tbody>
    </table>
    <ul class="memo-ul">
      <li><code>^</code> 仅在紧跟字符集开头的 <code>[</code> 时需转义。</li>
      <li><code>-</code> 仅在两个字母或两个数字之间时需转义。</li>
    </ul>

    <h2>量词</h2>
    <table class="memo-table">
      <thead><tr><th>表达式</th><th>说明</th></tr></thead>
      <tbody>
        <tr><td><code>{2}</code></td><td>恰好 2 次</td></tr>
        <tr><td><code>{2,}</code></td><td>至少 2 次</td></tr>
        <tr><td><code>{2,7}</code></td><td>至少 2 次但不超过 7 次</td></tr>
        <tr><td><code>*</code></td><td>0 次或更多</td></tr>
        <tr><td><code>+</code></td><td>1 次或更多</td></tr>
        <tr><td><code>?</code></td><td>恰好 0 次或 1 次</td></tr>
      </tbody>
    </table>
    <p>量词位于被修饰的表达式之后。</p>

    <h2>边界</h2>
    <table class="memo-table">
      <thead><tr><th>表达式</th><th>说明</th></tr></thead>
      <tbody>
        <tr><td><code>^</code></td><td>字符串开头</td></tr>
        <tr><td><code>$</code></td><td>字符串结尾</td></tr>
        <tr><td><code>\b</code></td><td>单词边界</td></tr>
      </tbody>
    </table>
    <p>单词边界的匹配规则：</p>
    <ul class="memo-ul">
      <li>字符串开头且首字符是 <code>\w</code>。</li>
      <li>相邻两字符之间，前一个是 <code>\w</code> 而后一个是 <code>\W</code>。</li>
      <li>字符串结尾且末字符是 <code>\w</code>。</li>
    </ul>

    <h2>匹配</h2>
    <table class="memo-table">
      <thead><tr><th>表达式</th><th>说明</th></tr></thead>
      <tbody>
        <tr><td><code>foo\|bar</code></td><td>匹配 <code>foo</code> 或 <code>bar</code></td></tr>
        <tr><td><code>foo(?=bar)</code></td><td>匹配后面紧跟 <code>bar</code> 的 <code>foo</code>（正向先行断言）</td></tr>
        <tr><td><code>foo(?!bar)</code></td><td>匹配后面不紧跟 <code>bar</code> 的 <code>foo</code>（负向先行断言）</td></tr>
        <tr><td><code>(?&lt;=bar)foo</code></td><td>匹配前面是 <code>bar</code> 的 <code>foo</code>（正向后行断言）</td></tr>
        <tr><td><code>(?&lt;!bar)foo</code></td><td>匹配前面不是 <code>bar</code> 的 <code>foo</code>（负向后行断言）</td></tr>
      </tbody>
    </table>

    <h2>分组与捕获</h2>
    <table class="memo-table">
      <thead><tr><th>表达式</th><th>说明</th></tr></thead>
      <tbody>
        <tr><td><code>(foo)</code></td><td>捕获组：匹配并捕获 <code>foo</code></td></tr>
        <tr><td><code>(?:foo)</code></td><td>非捕获组：匹配 <code>foo</code> 但不捕获</td></tr>
        <tr><td><code>(foo)bar\1</code></td><td><code>\1</code> 是对第 1 个捕获组的反向引用，匹配 <code>foobarfoo</code></td></tr>
      </tbody>
    </table>
    <p><code>\N</code> 是对第 N 个捕获组的反向引用，捕获组从 1 开始编号。Go 中使用命名组写法 <code>(?P&lt;name&gt;...)</code>。</p>
  </article>
</template>

<style scoped>
.memo {
  line-height: 1.6;
}

.memo h2 {
  margin: 24px 0 8px;
  font-size: 22px;
  font-weight: 600;
  color: v-bind('themeVars.textColor1');
}

.memo h2:first-child {
  margin-top: 0;
}

.memo h3 {
  margin: 16px 0 8px;
  font-size: 16px;
  font-weight: 600;
  color: v-bind('themeVars.textColor1');
}

.memo p {
  margin: 8px 0;
  color: v-bind('themeVars.textColor2');
}

.memo p code,
.memo li code,
.memo td code {
  font-family: 'Cascadia Code', Consolas, 'Courier New', monospace;
  font-size: 13px;
  background: v-bind('themeVars.cardColor');
  padding: 1px 5px;
  border-radius: 3px;
  color: v-bind('themeVars.textColor1');
}

.memo-table {
  width: 100%;
  border-collapse: collapse;
  margin: 10px 0 16px;
  font-size: 14px;
}

.memo-table th,
.memo-table td {
  border: 1px solid v-bind('themeVars.borderColor');
  padding: 6px 10px;
  text-align: left;
  color: v-bind('themeVars.textColor2');
}

.memo-table th {
  font-weight: 600;
  color: v-bind('themeVars.textColor1');
}

.memo-ul {
  margin: 8px 0 16px;
  padding-left: 22px;
  color: v-bind('themeVars.textColor2');
}

.memo-ul li {
  margin: 4px 0;
}
</style>
