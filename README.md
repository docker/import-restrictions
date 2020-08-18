# import-restrictions

Restrict imports in your go project.

This tools helps maintainers keep dependencies clean when they have multiple packages inside the same repo that shouldn't depend on each other.

## Usage

By default this tool will search for a configuration file named `import-restrictions.yaml` in the current directory, you can give it another configuration file with the flag `--configuration`.

```sh
$ import-restrictions
$ import-restrictions --configuration my-configuration-file.yaml
```

## Configuration

The configuration file is a yaml file with an array of objects that must contain:

- `dir` the base directory of the package
- `forbiddenImports` a list of packages that are forbidden for all the packages inside `dir`.

For example:

```yaml
- dir: ./cmd
  forbiddenImports:
    - bytes
    - github.com/account/repo
```

The configuration is pretty self-explanatory, with a configuration file like the one above all the packages inside the directory `cmd` cannot import any of the `bytes` and `github.com/account/repo` packages.
