```yml
name: Deploy Mail Cloud function

on:
  push:
    branches:
      - master
    paths:
      - ".github/workflows/mailfunction.yml"
      - "keystone-mail/**"

jobs:
  build:
    runs-on: ubuntu-latest

    strategy:
      matrix:
        node-version: [10.x]

    steps:
      - uses: actions/checkout@v1
      - name: Use Node.js ${{ matrix.node-version }}
        uses: actions/setup-node@v1
        with:
          node-version: ${{ matrix.node-version }}
          registry-url: https://registry.npmjs.org/
      - name: install global dependencies
        run: |
          npm add -g pnpm
          npm add -g @microsoft/rush
          npm install -g @keystone.sh/cli@latest
      - name: rush install, build
        # create an empty env to avoid warnings and process.exit = 1 on rush build
        run: |
          touch keystone-web/.env 
          rush install
      - name: deploy on Google
        run: |
          ks pull
          cd keystone-mail
          pnpm run deploy
        env:
          KEYSTONE_SHARED: ${{secrets.KEYSTONE_SHARED}}
```