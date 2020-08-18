# import-restrictions

Restrict imports in your go project.

## Usage

By default this tool will search for a configuration file named `import-restrictions.yaml` in the current directory, you can give it another configuration file with the flag `--configuration`.

```sh
$ import-restrictions
$ import-restrictions --configuration my-configuration-file.yaml
```

## Configuration

```yaml
- dir: ./cmd
  forbiddenImports:
    - bytes
    - github.com/account/repo
```

The configuration is pretty self-explanatory, with a configuration file like the one above all the packages inside the directory `cmd` cannot import any of the `bytes` and `github.com/account/repo` packages.
