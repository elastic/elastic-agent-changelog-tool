# Configuration options

`elastic-agent-changelog-tool` has configuration options available to change its behaviour.

All settings are managed via the [`settings`][settings] package, using [`spf13/viper`][viper].
Configurations are bound to environment variables with same name and `ELASTIC_AGENT_CHANGELOG` prefix using [`viper.BindEnv`][bindenv].

This CLI supports and adhere to cross platform XDG Standard provided by [`OpenPeeDeeP/xdg`][xdg].

| Settings | Default value | Description |
|---|---|---|
|`fragment_location`|`$GIT_REPO_ROOT/changelog/fragments`|The location of changelog fragments used by the CLI. By default `fragment_root` + `fragment_path`.|
|`fragment_path`|`changelog/fragments`|The path in `fragment_root` where to locate changelog fragments.|
|`fragment_root`|`$GIT_REPO_ROOT`|The root folder for `fragment_location`.|

## Configuration file

Add a `config.changelog.yaml` file to the repo where you're generating release notes.

When generating Markdown files, at a minimum you should set the following settings:

| Setting | Default value | Description |
|---|---|---|
| `owner` (required) | `elastic` | The owner of the GitHub repo. |
| `repo` (required) | ‒ | The name of the GitHub repo. |
| `rendered_changelog_destination` | `changelog` | The directory where you want to put the generated files.<br><br>When generating Markdown files, this should probably be `docs/release-notes/_snippets`. |

## Supported Environment Variables

`elastic-agent-changelog-tool` uses some environment variables that can be set.

|Name|Default|Note|
|---|---|---|
|`GIT_REPO_ROOT`|Git repository root folder|This value is computed at each execution to retrieve the repository root folder.|

[bindenv]: https://pkg.go.dev/github.com/spf13/viper#BindEnv
[settings]: ../internal/settings/settings.go
[xdg]: https://pkg.go.dev/github.com/OpenPeeDeeP/xdg
[viper]: https://pkg.go.dev/github.com/spf13/viper
