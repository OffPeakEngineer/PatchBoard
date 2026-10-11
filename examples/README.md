<!-- Copyright (c) 2026 Andrew David LeTourneau; MIT OR Zlib -->
# Publish your project's PatchBoard

Your project owns its `tasks/` directory. CI builds PatchBoard in a separate
checkout, runs its built-in HTML summary over **your** tasks, and saves a
single `public/index.html`. Publish that file to your project's Pages site or
any static web host. Visitors need no Node, JavaScript, or folder permissions.

| Example | Where to put it | What it does |
| --- | --- | --- |
| [GitHub Pages](github-pages.yml) | `.github/workflows/patchboard-pages.yml` | Builds on pushes and PRs; deploys the default branch to GitHub Pages |
| [GitLab Pages](gitlab-pages.yml) | `.gitlab/patchboard-pages.yml`, included by `.gitlab-ci.yml` | Builds branch/MR artifacts; deploys default-branch pushes to GitLab Pages |
| Local export | `patchboard summary --type=html --output public/index.html` | Produces one HTML file for local preview or another static host |

These examples only publish a board. Your project does not need npm package
metadata, Conventional Commits, semantic-release, or PatchBoard's own release
pipeline. CI uses Go to build the standalone binary in a separate checkout.
Local exports use an installed PatchBoard binary.

## Prepare your board

With the PatchBoard CLI installed, run `patchboard init` in your project if it
does not have a board yet. Commit the task tree, including `board.yml`, configured
state directories, Markdown tasks, and `kanban.html`. For an existing board,
`patchboard fix --dry-run` previews template updates and `patchboard fix` applies
them. Empty state directories need a tracked file such as `.gitkeep` to survive
checkout.

The exporter uses PatchBoard's Go task scanner and embedded board styles. It
does not require `kanban.html` to exist. The CI examples follow `main` by
default; set `PATCHBOARD_REF` to a commit or release containing `summary` for
repeatable builds. Earlier releases may predate the built-in exporter.

CI needs Git and Go 1.25 or newer to build PatchBoard. Board export needs no
Node or browser installation. All task text, including completed tasks, is
embedded in the HTML; configure Pages access for the intended audience.

## GitHub Pages

1. Copy `github-pages.yml` to `.github/workflows/patchboard-pages.yml`.
2. In **Settings → Pages → Build and deployment**, choose **GitHub Actions**.
3. Allow the `github-pages` environment to deploy your default branch.
4. Push the workflow and your task files. No release token or custom secret is
   needed: the build has read access and the deployment uses Pages/OIDC permissions.

The workflow detects your default branch. PRs and feature branches produce a
`github-pages` artifact without deploying. A manual run on the default branch
can redeploy. Change `TASK_DIRECTORY` for a board such as `planning/tasks`.
The second checkout is tooling only; the render command runs from your project
root and explicitly passes its task directory.

Use the site URL reported by the deployment, including the exact path casing.
The generated board's links are internal anchors, so repository subpaths need no
base-URL configuration. See the [GitHub Pages workflow guide](https://docs.github.com/en/pages/getting-started-with-github-pages/using-custom-workflows-with-github-pages).

## GitLab Pages

Copy `gitlab-pages.yml` to `.gitlab/patchboard-pages.yml`. For a project without
an existing pipeline, this `.gitlab-ci.yml` is sufficient:

```yaml
stages: [build, deploy]
include:
  - local: .gitlab/patchboard-pages.yml
```

For an existing pipeline, add the include and ensure its `stages` list contains
`build` and `deploy`, or change the example's stages to match yours. Existing
`workflow:rules` must permit branch pushes and merge requests. The jobs avoid
inheriting your application's `before_script`, caches, and artifacts. The
publish job waits for the board build; add your application's test jobs to its
`needs` list if they must also pass before deployment.

Enable Pages for the project. This example requires GitLab 17.9+ for nested
`pages.publish`, and uses the runner's built-in Pages publication without a
custom deployment token. Edit the build job's `TASK_DIRECTORY` for another
board path. MR/feature builds retain `public/index.html` as a downloadable
artifact for one week. Tag pipelines do not publish the board.
See [GitLab's Pages configuration](https://docs.gitlab.com/ci/yaml/#pagespublish).

## Local export or another static host

With PatchBoard installed, run from your project root:

```sh
patchboard summary --type=html --output public/index.html
# Custom task folder and output; paths with spaces work too:
patchboard -C planning/tasks summary --type=html --output site/board/index.html
```

Open the generated HTML directly or upload it to any static host. The published
output is only that HTML file; the tooling checkout does not belong in the site's
upload directory. See [Getting Started](../docs/getting-started.md) for binary
installation or source builds.

## Add the board to an existing website

A project normally has one Pages site. If your project already publishes docs
or an application, integrate the export into that site's existing build instead
of adding another workflow that replaces the site:

```sh
# First build your existing site into public/, then add the board:
patchboard summary --type=html --output public/board/index.html
# Your existing deployment now uploads the whole public/ directory once.
```

The tooling checkout and build steps come from either CI example. Link
from the site's home page to `board/`. Export after the site generator has
finished so it does not erase the board. Likewise, if your application stores
source files in `public/`, choose a dedicated output directory such as
`board-site/`, and change the artifact path and Pages publish path to match.

Add generated output and `.patchboard-tools/` to your project's `.gitignore`.
If you run `patchboard lint` after generating the site, add those directories to
`ignore_dirs` in `tasks/board.yml`, preserving your other exclusions. Refresh the
published snapshot with a new build whenever task files change.
