# Go Change Detection

## Description
- This is a CLI based tool that will determine which files changed between the current and last commit to main, or all the commits in a given pull request
- This will also detect all changed files in any given merge commit, as if they were all squashed

## Usage
- This is a standalone CLI, but was specifically developed to run in a Jenkins stage such as

``` groovy
stage('Change Detection - Release to Testing - Main') {
        environment {
            CI = "true"
        }
        options {
            skipDefaultCheckout(true)
        }
        steps {
          container("gather-changeset") {
            timeout(time: 10, unit: 'MINUTES') {
                sh "export DEBUG=true && gitActions app --repo-path=${workspace}/pac --branch=main"
                script {
                dirs = readFile("changeSet").readLines()
                }
                echo "Changed directories are: ${dirs}"
              }
            }
          }
          when {
            anyOf{
              branch 'main'
            }
          }
      }
```

- This will create a variable `dirs` that can be referenced in subsequent stages in a when condition, like so:

```groovy
when {
    expression {
      return dirs.contains("services/${SERVICE}".toString())
    }
}
```

- This works best in a matrix, where a list of services is a matrix axis:

```groovy
matrix {
  axes {
    axis {
      name 'SERVICE'
      values 'outsim', 'outsim-gateway'
    }
  }
  stages {...}
```

## Available Commands
``` console
gitActions --help
A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Usage:
gitActions [command]

Available Commands:
app         A brief description of your command
completion  Generate the autocompletion script for the specified shell
help        Help about any command
terraform   A brief description of your command

Flags:
-h, --help     help for gitActions
-t, --toggle   Help message for toggle

Use "gitActions [command] --help" for more information about a command.
```