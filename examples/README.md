<!-- Copyright (c) 2026 Andrew David LeTourneau; MIT OR Zlib -->
# Publish your project's PatchBoard

Your project owns its `tasks/` directory. CI checks out PatchBoard's exporter
separately, opens **your** `tasks/kanban.html`, reads **your** tasks, and saves a
single `public/index.html`. Publish that file to your project's Pages site or
any static web host. Visitors need no Node, JavaScript, or folder permissions.

| Example | Where to put it | What it does |
| --- | --- | --- |
| [GitHub Pages](github-pages.yml) | `.github/workflows/patchboard-pages.yml` | Builds on pushes and PRs; deploys the default branch to GitHub Pages |
| [GitLab Pages](gitlab-pages.yml) | `.gitlab/patchboard-pages.yml`, included by `.gitlab-ci.yml` | Builds branch/MR artifacts; deploys default-branch pushes to GitLab Pages |
| [Local export](export-local.sh) | Anywhere, such as `scripts/export-board.sh` | Produces one HTML file for local preview or another static host |

These examples only publish a board. Your project does not need Go, npm package
metadata, Conventional Commits, semantic-release, or PatchBoard's own release
pipeline. Node and Chromium are build tools installed in a separate checkout;
your project's package files are untouched.

## Prepare your board

With the PatchBoard CLI installed, run `patchboard init` in your project if it
does not have a board yet. Commit the task tree, including `board.yml`, configured
state directories, Markdown tasks, and `kanban.html`. For an existing board,
`patchboard fix --dry-run` previews template updates and `patchboard fix` applies
them. Empty state directories need a tracked file such as `.gitkeep` to survive
checkout.

The exporter uses the browser template's `loadHandle` function and DOM structure.
Use a compatible PatchBoard-generated `kanban.html`; an independently rewritten
board needs an adapted exporter. The examples pin the exporter to commit
`097f0d865eec8c7e0151efd270132b5cf4f95a06`, which includes the static export tool.
Old release tags may predate that tool. Upgrade the pin and the board template
together. Set `PATCHBOARD_REF` to a reviewed newer commit or release containing
`scripts/export-board.mjs` and the npm lockfile when upgrading.

The CI examples use current stable Node and actions with Node 24 runtimes. For
local exports, install Git and Node 24.10+ (current stable recommended). The
first build downloads npm dependencies and Chromium. All task text displayed
by the board, including completed tasks, is embedded in the HTML; configure
Pages access for the intended audience.

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

Copy `export-local.sh` into your project and run it from your project root:

```sh
sh scripts/export-board.sh
# Custom task folder and output; paths with spaces work too:
sh scripts/export-board.sh planning/tasks site/board/index.html
```

The script downloads the pinned tooling into a temporary checkout, installs its
locked dependencies and Chromium, exports your board, then removes that checkout.
Chromium and npm may retain their normal download caches. Linux machines need
Chromium's OS libraries; the CI examples install these with Playwright's
`--with-deps` option. The local script does not install system packages.

Open the generated HTML directly or upload it to any static host. The published
output is only that HTML file; neither the tooling checkout nor `node_modules`
belongs in the site's upload directory. Use `PATCHBOARD_SOURCE` to select your
own accessible Git mirror or an absolute local checkout path if needed.

## Add the board to an existing website

A project normally has one Pages site. If your project already publishes docs
or an application, integrate the export into that site's existing build instead
of adding another workflow that replaces the site:

```sh
# First build your existing site into public/, then add the board:
node .patchboard-tools/scripts/export-board.mjs tasks public/board/index.html
# Your existing deployment now uploads the whole public/ directory once.
```

The tooling checkout and installation steps come from either CI example. Link
from the site's home page to `board/`. Export after the site generator has
finished so it does not erase the board. Likewise, if your application stores
source files in `public/`, choose a dedicated output directory such as
`board-site/`, and change the artifact path and Pages publish path to match.

Add generated output and `.patchboard-tools/` to your project's `.gitignore`.
If you run `patchboard lint` after generating the site, add those directories to
`ignore_dirs` in `tasks/board.yml`, preserving your other exclusions. Refresh the
published snapshot with a new build whenever task files change.
