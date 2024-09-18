
set -e

#setup project folder
rm -rf ./app
mkdir -p ./app
cd ./app

#setup editor
echo """
root = true

[*]
charset = utf-8
indent_style = space
indent_size = 4
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true
quote_type = single
""" > .editorconfig

#setup .git
git init
touch .gitignore

#setup npm
npm init -y
mkdir -p ./src
echo build >> node_modules
npm install --save-dev @types/node

#setup typescript
npm install --save-dev typescript
npx tsc --init
sed -i 's|// "outDir": "./"| "outDir": "./build"|' ./tsconfig.json
sed -i 's/"compilerOptions": {/"include": [\n    ".\/src"\n  ],\n  "compilerOptions": {/' ./tsconfig.json
echo """export function add(a: number, b: number): number {
  return a + b
}

console.log(process.env['NODE_ENV'])
""" >> ./src/index.ts
sed -i 's/"main": "index.js"/"main": ".\/build\/index.js"/' ./package.json
sed -i 's/"test": .*/"build": "tsc --build"/' ./package.json
sed -i 's/"scripts": {/"scripts": {\n    "start": "node --env-file=.env .\/build\/index.js",/' ./package.json
echo build >> .gitignore
npm run build
npm run start

#setup typescript:watch
npm install --save-dev ts-node-dev
sed -i 's/"scripts": {/"scripts": {\n    "start:watch": "ts-node-dev --env-file=.env --respawn .\/src\/index.ts",/' ./package.json

#setup testing
npm install --save-dev @types/glob
mkdir -p ./tests
echo """import { describe, it } from 'node:test'
import assert from 'node:assert'
import { add } from '../src/index'

describe('add', () => {
    it('should detect static defined groups', () => {
        assert.strictEqual(add(1, 2), 3)
    })
})
""" > ./tests/index.test.ts
echo """
import process from 'node:process'
import {run} from 'node:test'
import {spec as SpecReporter} from 'node:test/reporters'
import {pipeline} from 'node:stream/promises'
import {glob} from 'glob'

async function main() {
  let fail = false

  const source = run({
      concurrency: true,
      files: glob.sync('tests/**/*.test.ts'),
  }).on('test:fail', () => {
        fail = true
  })

  const reporter = new SpecReporter()
  const destination = process.stdout
  await pipeline(source, reporter, destination)
  if (fail) throw new Error('Tests failed')
}

main()
""" >> ./tests/runner.ts
sed -i 's/"scripts": {/"scripts": {\n    "test": "ts-node-dev --env-file=.env .\/tests\/runner.ts",/' ./package.json
sed -i 's/"scripts": {/"scripts": {\n    "test:watch": "ts-node-dev --env-file=.env --respawn .\/tests\/runner.ts",/' ./package.json
npm run test

#setup formating
npm install --save-dev prettier
echo """{
  \"semi\": false,
  \"singleQuote\": true,
  \"tabWidth\": 2
}""" > .prettierrc.json
sed -i 's/"scripts": {/"scripts": {\n    "format": "prettier . --write",/' ./package.json
echo """node_modules
output""" > .prettierignore
npm run format

#setup linting
npm install --save-dev eslint @eslint/js globals typescript-eslint
echo """import globals from \"globals\"
import pluginJs from \"@eslint/js\"
import tseslint from \"typescript-eslint\"

export default [
  {files: [\"**/*.{ts}\"]},
  {languageOptions: { globals: globals.node }},
  pluginJs.configs.recommended,
  ...tseslint.configs.recommended,
]
""" > eslint.config.mjs
sed -i 's/"scripts": {/"scripts": {\n    "lint": "eslint .\/src",/' ./package.json
npm run lint

#setup env
echo "NODE_ENV=production" > .env
echo "NODE_ENV=development" > .env.dev
echo .env >> .gitignore
npm run build
npm run start
