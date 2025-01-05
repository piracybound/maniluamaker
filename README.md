# maniluamaker

maniluamaker is a simple go program for processing `.manifest` files and `config.vdf` files, generating a lua script that combines manifest and decryption key data.

## Build

- Ensure you have Go installed on your system. [Install Go](https://golang.org/doc/install) if you haven't already.

1. Clone the repository to your local machine.
2. Build the program:
   ```bash
   go build .
   ```

## Usage

1. Drag and drop `.manifest` files and a `config.vdf` file onto the program executable:
2. The program will:
   - Process the `.manifest` files to extract depot IDs and manifest IDs.
   - Process the `config.vdf` file to extract depot IDs and decryption keys.
   - Match depot IDs between the two and generate a lua script.
3. Enter the game's APP ID when prompted.
4. The output file will be named `<APP_ID>.lua` and saved in the same directory as the program.

## Example Input and Output

### Input Files

**File 1: `123456_789012.manifest`**
**File 2: `config.vdf`**

```plaintext
"123456" {
    "DecryptionKey" "ABCDEF1234567890"
}
```

### Generated Output File: `12345.lua`

```lua
-- manifest & lua provided by: https://www.piracybound.com/discord
-- via manilua
addappid(12345)
addappid(123456, 1, "ABCDEF1234567890")
setManifestid(123456, "789012", 0)
```

<sub>maniluamaker © 2025 by piracybound is licensed under <a href="https://github.com/piracybound/maniluamaker/blob/main/LICENSE">CC BY-ND 4.0</a></sub>
