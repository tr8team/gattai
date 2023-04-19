# Gattai

Popularized by Japanese giant robot-type shows, the term **_literally refers to the act of merging two separate entities._**

As the name suggested, this tool aims to simplify the process of pulling information from various sources and putting them together to form the configuration as well as sensitive information needed for deployment purposes.

**[To be added]: Sync secrets between your Kubernetes Pods and an External Secret Store (in our case it is AWS Secrets Manager)**

There are 2 types of files needed for Gattai to work

- `GattaiFile` defined the target to be generated/retrieved

  ```yaml
  version: v1

  temp_folder: ./localtmp

  enforce_targets:
    dev:
      - mysql_password
      - create_envfile

  repos:
    github-src:
      src: git
      config:
        url: https://github.com/tr8team/gattai.git
        branch: develop
        dir: gattai_actions

  targets:
    dev:
      mysql_password:
        action: github-src/op_cli_item_get
        vars:
          identifier: "Okta Local Testing"
          label: password
          vault: Private
      create_envfile:
        action: github-src/write_temp_file
        vars:
          filename: .env.remote-testing
          content: |
            DB_HOST=engineering-remote-testing-transactional-database
            DB_USERNAME=admin
            DB_PASSWORD={{ fetch .Targets.dev.mysql_password }}
            REDIS_HOST=engineering-remote-testing-cache-master
  ```

- `ActionFile` defines the steps needed by the user to retrieve the information needed

  ```yaml
  version: v1
  type: CommandLineInterface
  params:
    required:
      identifier:
        desc: The one password identifier
        type: string
      label:
        desc: The one password label
        type: string
      vault:
        desc: The one password vault
        type: string
  spec:
    test:
      expected:
        condition: equal
        value: 1
      cmds:
        - command: op
          {{- with .Vars }}
          args:
            - item list
            {{- with .identifier }}
            - "| grep \"{{ . }}\" | wc -l"
            {{- end }}
          {{- end }}
    exec:
      cmds:
        - command: op
          {{- with .Vars }}
          args:
            {{- with .identifier }}
            - item get "{{ . }}"
            {{- end }}
            {{- with .label }}
            - --fields label={{ . }}
            {{- end }}
            {{- with .vault }}
            - --vault {{ . }}
            {{- end }}
          {{- end }}
  ```

## Features

- **Easy to contribute and share** - create any git repository with `ActionFiles`, and you are ready to share your code with others
- **Enforce targets** - this feature allows a consistent implementation of targets across different deployments without the user forgetting to implement them.
- **Parameter enforcement** - this help to check whether any required arguments are missing and alert the user.
- **Testing and Validation** - Able to test and validate the existence of the source of information to ensure no issues are the result of missing information.
- **Inheritance** - Able to inherit from existing action files to further simplify the variables needed to trigger a source retrieval.

# Installation

## Nix Package

```nix
atomi = (
  with import (fetchTarball "https://github.com/kirinnee/test-nix-repo/archive/refs/tags/v10.0.0.tar.gz");
  {
    inherit
    gattai;
  }
);
```

# Usage

## Getting started

Create a file called `GattaiFile.yaml` or `<any name>.yaml` in the root of your project. The following contains the `minimum` configuration needed to get started.

```yaml
version: v1

repos:
  github-src:
    src: git
    config:
      url: https://github.com/tr8team/gattai.git
      branch: develop
      dir: gattai_actions

targets:
  dev:
    mysql_password:
      action: github-src/op_cli_item_get
      vars:
        identifier: "Okta Local Testing"
        label: password
        vault: Private
```

**Minimum Requirement**

- `version` This is the gattai file format version
- `repos` This is the repo to the action files, currently only `local` and `git` src are supported
- `targets` This is the list of targets that the user can specify to trigger, retrieve and construct the final outcome. Running gattai is as simple as the following:

```bash
gattai run dev mysql_password <any name>.yaml
```

**Additional Requirement**

- `temp_folder` This is the folder where the temporary generated file will be stored
- `enforce_targets` This is the list of targets being enforced and are mandatory to be implemented or else the gattai application will make noise

## Supported file names

Gattai will look for the following filename by default if no filename is provided

- `GattaiFile.yaml`

Else any filename ending with the `<any name>.yaml` extension is supported

## How to add target

```yaml
targets:
  <namespace>:
    <target>:
      action: github-src/op_cli_item_get
      vars:
        identifier: "Okta Local Testing"
        label: password
        vault: Private
```

- `action` Name of the action file in the following format `<repo>/<path_to_action_file_wo_ext>`
  - `repo`: Using the example above, repo is `github-src`
  - `path_to_action_file_wo_ext`: Using the example above, the action file without extension is `op_cli_item_get`
- `vars` These are the argument needed by the action file to perform its yaml. The variables required can be found in the `params` section in the action file.

## Action Files

Currently there are 2 type of action file, they are `CommandLineInterface` and `DerivedInterface`.

### CommandLineInterface

The `CommandLineInterface` uses [this shell interpreter](https://github.com/mvdan/sh) to execute command line with support for shell operations. Along with the help from Go Template, the user can customised their command line before they trigger it.

```yaml
version: v1
type: CommandLineInterface
params:
  required:
    region:
      desc: The regiion where the resource is at
      type: string
    identifier:
      desc: The rds cluster identifier
      type: string
  optional:
    query:
      desc: The rds cluster query
      type: string
spec:
  test:
    expected:
      condition: equal
      value: 1
    cmds:
      - command: aws
        {{- with .Vars }}
        args:
          {{- with .region }}
          - --region {{ . }}
          {{- end }}
          - rds describe-db-clusters
          {{- with .identifier }}
          - --db-cluster-identifier {{ . }}
          {{- end }}
          - "| jq '. | length'"
        {{- end }}
  exec:
    cmds:
      - command: aws
        {{- with .Vars }}
        args:
          {{- with .region }}
          - --region {{ . }}
          {{- end }}
          - rds describe-db-clusters
          {{- with .identifier }}
          - --db-cluster-identifier {{ . }}
          {{- end }}
          {{- with .query }}
          - --query "DBClusters[0].{{ . }}"
          {{- end }}
        {{- end }}
```

- `version` The version of action file
- `type` The type of action file
- `params` The list of parameters to be provided in order to use the action file
  - `required` These are the required fields to use the action file
  - `optional` These are the optional fields to use the action file. However, more often than not, the fields are meant for the user to `choose at least one` of the settings rather than totally `optional`.
- `spec` These are the specification use for this action file

  - `test` This is where the logic for verifying the existence of the information source via running the `validate` sub-command

    - `condition` Currently 6 comparison condition are supported:

      | equal         | Check for result equal to the value                |
      | ------------- | -------------------------------------------------- |
      | not equal     | Check for result not equal to the value            |
      | contain       | Check for sub-string value exist in the result     |
      | not contain   | Check for sub-string value not exist in the result |
      | int less than | Check for result is integer less than the value    |
      | int more than | Check for result is integer more than the value    |

    - `value` The value which the result will be compared to

  - `exec` This is where the logic for retrieving the information source via running the `run` sub-command

### DerivedInterface

The `DerivedInterface` wrap another interface such as `CommandLineInterface` and simplify the parameters required by allowing user to hard-code certain arguments.

```yaml
version: v1
type: DerivedInterface
params:
  required:
    region:
      desc: The region where the resource is at
      type: string
    identifier:
      desc: The content to be save if any
      type: string
spec:
  repo:
    src: local
    config:
      dir: ./
  override_test:
    expected:
      condition: equal
      value: 1
    cmds:
      - command: aws
        {{- with .Vars }}
        args:
          {{- with .region }}
          - --region {{ . }}
          {{- end }}
          - rds describe-db-clusters
          {{- with .identifier }}
          - --db-cluster-identifier {{ . }}
          {{- end }}
          - "| jq '. | length'"
        {{- end }}
  inherit_exec:
    action: aws/query/rds_cluster
    {{- with .Vars }}
    vars:
      region: "{{ .region }}"
      identifier: "{{ .identifier }}"
      query: Endpoint
    {{- end }}
```

- `version` The version of action file
- `type` The type of action file
- `params` The list of parameters to be provided in order to use the action file
  - `required` These are the required fields to use the action file
  - `optional` These are the optional fields to use the action file. However, more often than not, the fields are meant for the user to `choose at least one` of the settings rather than totally `optional`.
- `spec` These are the specification use for this action file
  - `repo` is where the action file you want to derived from
  - `override_test` is similar as `test`, just that instead of running the base interface test, it run this test instead
  - `inherit_exec` is similar as `target` and user fill it up as if they are triggering another action file

# API References

Gattai command line tool have the following syntax

```bash
gattai run <namespace> <target> [Gattaifile.yaml]
gattai validate <namespace> <target> [Gattaifile.yaml]
```

## Run

This sub-command execute the target as specified

| Short | Flag        | Type | Default | Description                                      |
| ----- | ----------- | ---- | ------- | ------------------------------------------------ |
| -e    | --enforce   | bool | false   | Flag to enforce targets when running the targets |
| -k    | --keep-temp | bool | false   | Flag to keep temporary generated file            |

## Validate

This sub-command check the target to see whether they existing and check all enforcement by default

## Targets

Each target is separated by the following

- `<namespace>` this tag help to group similar target under the same scope, you can use `all` to run all namespaces
- `<target>` this tag specify a target, you can use `all` to run all targets
