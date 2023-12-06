#!/bin/bash

set -e

git init

PACKAGE_JSON="package.json"
TSCONFIG_JSON="tsconfig.json"
GITIGNORE=".gitignore"
DOT_ENV=".env"
BUILD="build"
NODE_MODULES="node_modules"
ESLINTRC=".eslintrc.json"
ESLINTIGNORE=".eslintignore"
PRETTIERRC=".prettierrc.json"
PRETTIERIGNORE=".prettierignore"

echo "Create gitignore"
rm -rf "${GITIGNORE}"
echo "${BUILD}" >> "${GITIGNORE}"
echo "${NODE_MODULES}" >> "${GITIGNORE}"
echo "${DOT_ENV}" >> "${GITIGNORE}"

echo "Create env"
rm -rf "${DOT_ENV}"
rm -rf "${DOT_ENV}.sample"
echo "NODE_ENV=development" >> "${DOT_ENV}"
echo "NODE_ENV=development" >> "${DOT_ENV}.sample"

echo "Configure eslint"
rm -rf "${ESLINTIGNORE}"
echo "${NODE_MODULES}" >> "${ESLINTIGNORE}"
echo "${BUILD}" >> "${ESLINTIGNORE}"

rm -rf "${ESLINTRC}"
echo "{
  \"root\": true,
  \"parser\": \"@typescript-eslint/parser\",
  \"plugins\": [
    \"@typescript-eslint\"
  ],
  \"extends\": [
    \"eslint:recommended\",
    \"plugin:@typescript-eslint/eslint-recommended\",
    \"plugin:@typescript-eslint/recommended\"
  ],
  \"rules\": {
    \"@typescript-eslint/no-explicit-any\": 0,
    \"@typescript-eslint/no-inferrable-types\": 0
  }
}" >> "${ESLINTRC}"

echo "Configure prettier"
rm -rf "${PRETTIERIGNORE}"
echo "${NODE_MODULES}" >> "${PRETTIERIGNORE}"
echo "${BUILD}" >> "${PRETTIERIGNORE}"

rm -rf "${PRETTIERRC}"
echo "{
  \"semi\": false,
  \"singleQuote\": true,
  \"tabWidth\": 2
}" >> "${PRETTIERRC}"

echo "Regenerate package json"
rm -rf "${NODE_MODULES}"
rm -rf package-lock.json
rm -rf "${PACKAGE_JSON}"
npm init -y
npm install --save-dev ts-node-dev typescript eslint prettier @typescript-eslint/eslint-plugin

echo 'Parse package json'
cat "${PACKAGE_JSON}" | egrep -v "\"test\"" > ".${PACKAGE_JSON}"
cat ".${PACKAGE_JSON}" > "${PACKAGE_JSON}"
rm ".${PACKAGE_JSON}"
sed -i 's/"scripts": {/"scripts": {\n    "build": "tsc"/' "${PACKAGE_JSON}"
sed -i 's/"scripts": {/"scripts": {\n    "start:watch": "ts-node-dev --respawn .\/src\/index.ts",/' "${PACKAGE_JSON}"
sed -i 's/"scripts": {/"scripts": {\n    "start": "node .\/build\/src\/index.js",/' "${PACKAGE_JSON}"
sed -i 's/"scripts": {/"scripts": {\n    "format": "prettier . --write",/' "${PACKAGE_JSON}"
sed -i 's/"scripts": {/"scripts": {\n    "lint": "eslint . --ext .ts",/' "${PACKAGE_JSON}"
sed -i 's/"scripts": {/"scripts": {\n    "test:watch": "ts-node-dev --respawn .\/tests\/index.test.ts",/' "${PACKAGE_JSON}"
sed -i 's/"scripts": {/"scripts": {\n    "test": "node --test .\/build",/' "${PACKAGE_JSON}"

echo 'Regenerate tsconfig json'
rm -rf "${TSCONFIG_JSON}"
npx tsc --init

echo 'Parse tsconfig json'
sed -i 's/\/\/ \"outDir\": \".\/\"/\"outDir\": \".\/build\"/' "${TSCONFIG_JSON}"
sed -i 's/\/\/ \"declaration\": true/\"declaration\": true/' "${TSCONFIG_JSON}"

echo 'Generate code'
rm -rf ./src
mkdir ./src
echo "async function main(): Promise<void> {
    console.log('Hello world')
}

main().catch(console.error)
" > ./src/index.ts

rm -rf ./tests
mkdir ./tests
echo "import { describe, it } from 'node:test'
import assert from 'node:assert'

describe('module', () => {
    it('should be equals', () => {
        assert.equal(1, 1)
    })
    it('should not be equals', () => {
        assert.notEqual(1, 0)
    })
})
" > ./tests/index.test.ts
