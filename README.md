# bitcoin-wallet

Another way to back up your addresses with helpful data such as derivation path, and the address generation type from
your bitcoin wallet. This approach require **bitcoind** running and **bitcoin-cli** to catch the list of descriptors and have
legal access to it as the prerequisite.

## Steps

1. Download descriptors with public key and store it to file `descriptor_pubs.json` and copy it to directory `_data`
   under this repository:

```bash
./bitcoin-cli listdescriptors > descriptor_pubs.json
```

2. Download descriptors with root key and store it to file `descriptors.json` and copy it to directory `_data` under
   this repository:

```bash
# 1. unlock wallet
./bitcoin-cli walletpassphrase "your-wallet-passphrase" 60

# 2. download descriptors
./bitcoin-cli listdescriptors true > descriptors.json
```

3. Configure `_data/env.json`  with the correct value. You can provide it from `bitcoin.conf` file where your bitcoin
   base directory is.

4. Build project via command `make build-{your-os-version}`. You can find binary result on `_build` directory.

## Finally

After [steps](#steps) followed correctly and the structure should be:

```bash
.
├── _build
│   ├── btcd-darwin
│   └── btcd-linux
├── _data
│   ├── descriptor_pubs.json
│   ├── descriptors.json
│   └── env.json
```

Run binary from the root `_build/btcd-{os-version}` and then you can check result in `_data/derivation.csv` file.