<h2 align="center">license-tool</h2>

[中文](README_zh-CN.md)

## Project Introduction
1. The original intention of this project is to provide a `license authorization` tool for `delivery` programs
2. The project is divided into server and client. The client will generate a string of `signature codes`, get the server `request to generate a license license file and confuse the content`, and then the client will `de-obfuscate the license file and Verify file internal feature code`

## Development progress
- [x] server
- [x] client
- [x] Client `signature code` generation
- [x] server `license file` generated
- [x] The `license file` generated by the server is confusing
- [x] Client deobfuscates `license file`
- [x] The client verifies the signature inside the `license file`
- [ ] The client verifies the `license file` time
- [ ] The client verifies the remaining days of `license file` up to the current date
- [ ] logging
- [ ] The UTC time generated in `license` on the server side is converted to local time
- [ ] The server verifies whether `license file` is valid
- [ ] The server checks the `license permission` list
- [x] The server `license permission information` is stored in the database
- [ ] The client package is so and dll
- [ ] Request API doc
- [ ] Client Function API doc

## Installation start
```bash
# clone the project
https://github.com/keington/license-tool.git

# enter the project directory
cd license-tool/server

# install dependencies
go mod tidy

# develop
go run .
```

## Acknowledgments
Thanks to [JetBrain](https://www.jetbrains.com/) for the JetBrain Family Bucket Authorization License.

## Link me
If you have any questions, please ask questions in issues, and we will answer them regularly, or you can also send [mail](mailto:keington@outlook.com)

## License and Copyright
[MIT License](https://github.com/keington/license-tool/blob/cc897613c01f6ff7d2745ae1eb7303ff15a59d1c/LICENSE)

Copyright (c) 2023 Huaian Xu
