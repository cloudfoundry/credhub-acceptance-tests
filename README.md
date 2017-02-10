# Acceptance tests for CredHub

### Build the CLI first:

```
cd ../credhub-cli
make
```

### Run Tests locally

Target your local API by running:

```sh
./target_local.sh
```

If you want to target a different API you can edit the generated `config.json` file.

Runs local CredHub testing via:

```sh
./run_tests.sh
```
