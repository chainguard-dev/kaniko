# Kaniko Release Process

This document explains the Kaniko release process.

## Self-Serve Kaniko Release

Only members of the `chainguard-dev` organization and maintainers of the
`kaniko` repository can [create an official release of Kaniko](#official-release).

Other users must use the following steps:

1. Follow the [setup instruction to fork kaniko repository](https://github.com/chainguard-dev/kaniko/blob/master/DEVELOPMENT.md#getting-started)
2. Run the following `make` commands to build and push Kaniko image to your organization image repository and registry..
  ```shell
   REGISTRY=/YOUR-PROJECT make images
   ```
  The above command will build and push all the 3 kaniko images
  * YOUR-REGISTRY/YOUR-PROJECT/executor:latest
  * YOUR-REGISTRY/YOUR-PROJECT/executor:debug
  * YOUR-REGISTRY/YOUR-PROJECT/warmer:latest

3. You can choose to tag these images using `docker tag`. For example to tag
   `YOUR-REGISTRY/YOUR-PROJECT/executor:latest` as
   `YOUR-REGISTRY/YOUR-PROJECT/executor:v1.25.0self-serve`:

   ```shell
    docker tag YOUR-REGISTRY/YOUR-PROJECT/executor:latest YOUR-REGISTRY/YOUR-PROJECT/executor:v1.25.0self-serve
   ```
   
Change all usages of the official image to `OUR-REGISTRY/YOUR-PROJECT/executor:latest` for executor image and so on.
4. Finally, push your tagged images via docker. You could also use the Makefile target `push` to push these images like this
  ```shell
   REGISTRY=YOUR-REGISTRY/YOUR-PROJECT make images
  ```

<a id="official-release"></a>

## Kaniko Release Process

### Getting write access to Kaniko

In order to kick off kaniko release, you need to write access to Kaniko repository.

Once you have the correct access, you can kick off a release.


### Kicking of a release.

1. Create a release PR and update Changelog.

    In order to release a new version of Kaniko, you will need to first

    a. Create a new branch and bump the kaniko version in [Makefile](https://github.com/chainguard-dev/kaniko/blob/master/Makefile#L16)


    In most cases, you will need to bump the `VERSION_MINOR` number.
    In case you are doing a patch release for a hot fix, bump the `VERSION_BUILD` number.

    b. Run the [script](https://github.com/chainguard-dev/kaniko/blob/master/hack/release.sh) to create release notes.
    ```
    ./hack/release.sh
    Collecting pull request that were merged since the last release: v1.25.0 (2020-08-18 02:53:46 +0000 UTC)
    * ...
    Huge thank you for this release towards our contributors: 
    - Contributor 1
    - Contributor 2
    - ...
    ```
    Copy the release notes and update the [CHANGELOG.md](https://github.com/chainguard-dev/kaniko/blob/master/CHANGELOG.md) at the root of the repository.

    c. Create a pull request and get it approved from Kaniko maintainers.

2. Once the PR is approved and merged, create a release tag with name `vX.Y.Z` where
    ```
    X corresponds to VERSION_MAJOR
    Y corresponds to VERSION_MINOR
    Z corresponds to VERSION_BUILD
    ```
    E.g. to release 1.2.0 version of kaniko, please create a tag v1.2.0 like this
    ```
    git pull remote master
    git tag v1.2.0
    git push remote v1.2.0
    ```
3.  Pushing a tag to remote with above naming convention will trigger the Github workflow action defined [here](https://github.com/chainguard-dev/kaniko/blob/main/.github/workflows/images.yaml) It takes 20-30 mins for the job to finish and push images for Chainguard users.

```
cgr.dev/YOUR-ORG/kaniko:latest
gcr.io/kaniko-project/executor:latest
cgr.dev/YOUR-ORG/kaniko:vX.Y.Z
cgr.dev/YOUR-ORG/kaniko:debug
cgr.dev/YOUR-ORG/kaniko:debug-vX.Y.Z
cgr.dev/YOUR-ORG/kaniko-warmer
cgr.dev/YOUR-ORG/kaniko-warmer-vX.Y.Z
```

You can verify if the images are published using the `docker pull` command
```
docker pull cgr.dev/YOUR-ORG/kaniko:vX.Y.Z
docker pull cgr.dev/YOUR-ORG/kaniko-warmer:vX.Y.Z
```
In case the images are still not published, ping one of the kaniko maintainers and they will provide the cloud build trigger logs.

4. Finally, once the images are published, create a release for the newly created [tag](https://github.com/chainguard-dev/kaniko/tags) and publish it. 
Summarize the change log to mention, 
- new features added if any
- bug fixes, 
- refactors and 
- documentation changes
