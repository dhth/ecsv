# ecsv

[![Build Workflow Status](https://img.shields.io/github/actions/workflow/status/dhth/ecsv/build.yml?style=flat-square)](https://github.com/dhth/ecsv/actions/workflows/build.yml)
[![Vulncheck Workflow Status](https://img.shields.io/github/actions/workflow/status/dhth/ecsv/vulncheck.yml?style=flat-square&label=vulncheck)](https://github.com/dhth/ecsv/actions/workflows/vulncheck.yml)
[![Latest Release](https://img.shields.io/github/release/dhth/ecsv.svg?style=flat-square)](https://github.com/dhth/ecsv/releases/latest)
[![Commits Since Latest Release](https://img.shields.io/github/commits-since/dhth/ecsv/latest?style=flat-square)](https://github.com/dhth/ecsv/releases)

`ecsv` helps you quickly check the versions of your systems running in ECS tasks
across various environments.

![ecsv-terminal](https://github.com/user-attachments/assets/9faec97f-dda7-442c-a890-6059492b848b)

💾 Installation
---

**homebrew**:

```sh
brew install dhth/tap/ecsv
```

**go**:

```sh
go install github.com/dhth/ecsv@latest
```

⚡️ Usage
---

Create a configuration file that looks like the following.

```yaml
env-sequence: ["qa", "staging"]
systems:
- key: service-a
  envs:
  - name: qa
    aws-config-source: profile:::qa
    aws-region: eu-central-1
    cluster: 1brd-qa
    service: service-a-fargate
    container-name: service-a-qa-Service
  - name: staging
    aws-profile: qa
    aws-config-source: profile:::staging
    aws-region: eu-central-1
    cluster: 1brd-staging
    service: service-a-fargate
    container-name: service-a-staging-Service
- key: service-b
  envs:
  - name: qa
    aws-config-source: profile:::qa
    aws-region: eu-central-1
    cluster: 1brd-qa
    service: service-b-fargate
    container-name: service-b-qa-Service
  - name: staging
    aws-config-source: profile:::staging
    aws-region: eu-central-1
    cluster: 1brd-staging
    service: service-b-fargate
    container-name: service-b-staging-Service
```

By default, `ecsv` will try to find the config file at `~/.config/ecsv.yml`.

🔠 Output Formats
---

Besides the default ANSI output, `ecsv` can also output data in plaintext and
HTML formats.

```bash
ecsv -f table
```

![ecsv-table](https://github.com/user-attachments/assets/9003ab4c-09c0-44f8-b6a6-6933a0088f6a)

```bash
ecsv -f html > output.html
```

![ecsv-terminal](https://github.com/user-attachments/assets/dbde169a-3253-42cd-b5ff-0f2f99cecf58)

Read more about outputting HTML in the [examples](./examples/html-template)
directory.
