# v0.4.0 (2025-05-31)

## What's Changed

**Full Changelog**: https://github.com/gojue/moling/compare/v0.3.2...v0.4.0
<hr>

# v0.3.2 (2025-05-20)

## What's Changed

* feat: update command arguments and add new client configurations for Trae and Trae CN by @cfc4n
  in https://github.com/gojue/moling/pull/35
* feat: update command arguments and improve logging for service initialization by @cfc4n
  in https://github.com/gojue/moling/pull/37


**Full Changelog**: https://github.com/gojue/moling/compare/v0.3.1...v0.3.2
<hr>

# v0.3.1 (2025-04-22)

## What's Changed

* fix: handle parent process exit to ensure MCP Server shutdown by @cfc4n
  in [#33](https://github.com/gojue/moling/pull/33)
* feat: update default prompts for browser, command, and filesystem modules by @cfc4n
  in [#34](https://github.com/gojue/moling/pull/34)

**Full Changelog**: https://github.com/gojue/moling/compare/v0.3.0...v0.3.1
<hr>

# v0.3.0 (2025-04-18)

## What's Changed

* feat: add client configuration paths for Trae and update README by @cfc4n
  in [#27](https://github.com/gojue/moling/pull/27)
* feat: support Cursor IDE (MCP Client) by @cfc4n in [#28](https://github.com/gojue/moling/pull/28)
* feat: Improves code consistency, adds PID management, and enhances user prompts by @cfc4n
  in [#29](https://github.com/gojue/moling/pull/29)

**Full Changelog**: https://github.com/gojue/moling/compare/v0.2.0...v0.3.0
<hr>

# v0.2.0 (2025-04-13)

## What's Changed

* feat: add support for advanced logging configuration by @cfc4n in [#25](https://github.com/gojue/moling/pull/25)
* fix: resolve issue with incorrect breakpoint handling by @cfc4n in [#26](https://github.com/gojue/moling/pull/26)
* docs: update README with debugging examples by @cfc4n in [#27](https://github.com/gojue/moling/pull/27)

**Full Changelog**: https://github.com/gojue/moling/compare/v0.1.1...v0.2.0
<hr>

# v0.1.1 (2025-04-08)
## What's Changed

* feat: update CLI descriptions for clarity and accuracy by @cfc4n in https://github.com/gojue/moling/pull/17
* refactor: change log levels from Info to Debug and Warn for improved … by @cfc4n
  in https://github.com/gojue/moling/pull/19
* fix: update README files to replace image links with direct URLs by @cfc4n in https://github.com/gojue/moling/pull/20
* feat: add MCP Client configuration to install script by @cfc4n in https://github.com/gojue/moling/pull/22

**Full Changelog**: https://github.com/gojue/moling/compare/v0.1.0...v0.1.1
<hr>

# v0.1.0 (2025-04-04)

## What's Changed

* Improve builder by @cfc4n in https://github.com/gojue/moling/pull/16

**Full Changelog**: https://github.com/gojue/moling/compare/v0.0.8...v0.1.0
<hr>

# v0.0.8 (2025-04-04)

## What's Changed

* docs: add Japanese README file by @eltociear in https://github.com/gojue/moling/pull/13
* feat: implement log file rotation and enhance logging configuration by @cfc4n
  in https://github.com/gojue/moling/pull/14
* chore: update dependencies for chromedp and mcp-go to latest versions by @cfc4n
  in https://github.com/gojue/moling/pull/15

## New Contributors

* @eltociear made their first contribution in https://github.com/gojue/moling/pull/13

**Full Changelog**: https://github.com/gojue/moling/compare/v0.0.7...v0.0.8
<hr>

# v0.0.7 (2025-03-31)

## What's Changed

* feat: add logging enhancements and listen address flag for SSE mode by @cfc4n
  in https://github.com/gojue/moling/pull/11
* feat: add client command for automated MCP Server configuration management by @cfc4n
  in https://github.com/gojue/moling/pull/12

**Full Changelog**: https://github.com/gojue/moling/compare/v0.0.6...v0.0.7
<hr>

# v0.0.6 (2025-03-30)

## What's Changed

* feat: enhance service initialization and configuration loading by @cfc4n in https://github.com/gojue/moling/pull/9
* feat: update installation script and README for improved user experie… by @cfc4n
  in https://github.com/gojue/moling/pull/10

**Full Changelog**: https://github.com/gojue/moling/compare/v0.0.5...v0.0.6
<hr>

# v0.0.5 (2025-03-29)
## What's Changed

* fix #4 invalid data path by @cfc4n in https://github.com/gojue/moling/pull/5
* feat: add automated release notes generation in release workflow by @cfc4n in https://github.com/gojue/moling/pull/6
* Improve logging & error handling, rename configs, simplify service logic by @cfc4n
  in https://github.com/gojue/moling/pull/8

**Full Changelog**: https://github.com/gojue/moling/compare/v0.0.3...v0.0.5
<hr>

# v0.0.3 (2025-03-28)
## What's Changed

* fix: build failed with GOOS=windows by @cfc4n in https://github.com/gojue/moling/pull/1
* feat: refactor server initialization and directory creation logic by @cfc4n in https://github.com/gojue/moling/pull/2
* feat: enhance configuration management and logging in the MoLing server by @cfc4n
  in https://github.com/gojue/moling/pull/3

## New Contributors

* @cfc4n made their first contribution in https://github.com/gojue/moling/pull/1

**Full Changelog**: https://github.com/gojue/moling/compare/v0.0.2...v0.0.3
<hr>

# v0.0.2 (2025-03-27)

## What's Changed
- Enhanced command validation with global configuration and support for piped commands
- Added installation scripts for Windows/Linux and updated documentation
- Implemented directory creation checks for base path and subdirectories
- Corrected spelling of "MoLing" and "macOS" in documentation
- Updated filename format for saved images
- Refactored configuration management to use BasePath and improved browser initialization
- Added listen address configuration and improved server initialization
- Enhanced configuration management and logging
- Added configuration management and service registration
- Added browser service with logging
- Updated GitHub links in README for accurate repository references

**Full Changelog**: https://github.com/gojue/moling/compare/v0.0.1...v0.0.2
<hr>

# v0.0.1 (2025-03-23)

## What's Changed
1. support `filesystem` service
2. support `command line` service
3. add GitHub action workflow

