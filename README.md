# Bloodcoin

## How to start

```shell
git submodule init
git submodule update --remote
cd web
npm install
```

## How to

### List blocks

```shell
./bloodcoin-cli dump
```

### Add a sample prescription

```shell
./bloodcoin-cli prescription -data "$(cat sample/valid_prescription.json)"
```