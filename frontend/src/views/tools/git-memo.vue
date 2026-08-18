<script setup lang="ts">
import { useThemeVars } from 'naive-ui'

const themeVars = useThemeVars()
</script>

<template>
  <article class="memo">
      <h2>Configuration</h2>
      <p>Set the global config</p>
      <pre><code>git config --global user.name "[name]"
git config --global user.email "[email]"</code></pre>

      <h2>Get started</h2>
      <p>Create a git repository</p>
      <pre><code>git init</code></pre>
      <p>Clone an existing git repository</p>
      <pre><code>git clone [url]</code></pre>

      <h2>Commit</h2>
      <p>Commit all tracked changes</p>
      <pre><code>git commit -am "[commit message]"</code></pre>
      <p>Add new modifications to the last commit</p>
      <pre><code>git commit --amend --no-edit</code></pre>

      <h2>I’ve made a mistake</h2>
      <p>Change last commit message</p>
      <pre><code>git commit --amend</code></pre>
      <p>Undo most recent commit and keep changes</p>
      <pre><code>git reset HEAD~1</code></pre>
      <p>Undo the <code>N</code> most recent commit and keep changes</p>
      <pre><code>git reset HEAD~N</code></pre>
      <p>Undo most recent commit and get rid of changes</p>
      <pre><code>git reset HEAD~1 --hard</code></pre>
      <p>Reset branch to remote state</p>
      <pre><code>git fetch origin
git reset --hard origin/[branch-name]</code></pre>

      <h2>Miscellaneous</h2>
      <p>Renaming the local master branch to main</p>
      <pre><code>git branch -m master main</code></pre>
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

.memo p {
  margin: 8px 0;
  color: v-bind('themeVars.textColor2');
}

.memo p code {
  font-family: 'Cascadia Code', monospace;
}

.memo pre {
  margin: 10px 0 16px;
  padding: 15px 22px;
  overflow-x: auto;
  border-radius: 4px;
  background: v-bind('themeVars.cardColor');
  color: v-bind('themeVars.textColor1');
}

.memo pre code {
  font-family: 'Cascadia Code', monospace;
  white-space: pre;
}
</style>
