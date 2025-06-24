// postcss.config.cjs

const tailwindcss = require('@tailwindcss/postcss')
const autoprefixer = require('autoprefixer')

module.exports = {
  plugins: [
    tailwindcss(),    // ← 使用 @tailwindcss/postcss 提供的插件
    autoprefixer(),   // ← Autoprefixer 还是这样用
  ],
}

