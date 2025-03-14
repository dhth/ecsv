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

Or get the binaries directly from a
[release](https://github.com/dhth/ecsv/releases). Read more about verifying the
authenticity of released artifacts [here](#-verifying-release-artifacts).

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

🔐 Verifying release artifacts
---

In case you get the `ecsv` binary directly from a [release][4], you may want to
verify its authenticity. Checksums are applied to all released artifacts, and
the resulting checksum file is signed using
[cosign](https://docs.sigstore.dev/cosign/installation/).

Steps to verify (replace `A.B.C` in the commands listed below with the version
you want):

1. Download the following files from the release:

    - ecsv_A.B.C_checksums.txt
    - ecsv_A.B.C_checksums.txt.pem
    - ecsv_A.B.C_checksums.txt.sig

2. Verify the signature:

   ```shell
   cosign verify-blob ecsv_A.B.C_checksums.txt \
       --certificate ecsv_A.B.C_checksums.txt.pem \
       --signature ecsv_A.B.C_checksums.txt.sig \
       --certificate-identity-regexp 'https://github\.com/dhth/ecsv/\.github/workflows/.+' \
       --certificate-oidc-issuer "https://token.actions.githubusercontent.com"
   ```

3. Download the compressed archive you want, and validate its checksum:

   ```shell
   curl -sSLO https://github.com/dhth/ecsv/releases/download/vA.B.C/ecsv_A.B.C_linux_amd64.tar.gz
   sha256sum --ignore-missing -c ecsv_A.B.C_checksums.txt
   ```

3. If checksum validation goes through, uncompress the archive:

   ```shell
   tar -xzf ecsv_A.B.C_linux_amd64.tar.gz
   ./ecsv
   # profit!
   ```

≈ Related tools
---

- [ecscope](https://github.com/dhth/ecscope) lets you monitor ECS resources and
  deployments.
