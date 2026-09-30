# Releasing

Built by GitHub Actions (`.github/workflows/release.yml`). Every target builds
natively on its own runner - nothing cross-compiles, because Fyne needs CGO and a
matching C toolchain. For local builds see `BUILD.md`.

## Cutting a release

```bash
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3
```

The tag builds every enabled target, then publishes a GitHub Release with the
artifacts attached and auto-generated notes. Archives are named
`Cornealius-Eyeworth-<tag>-<os>-<arch>`:

| Platform | Runner | Archive | Contains | Built with |
|----------|--------|---------|----------|------------|
| `linux-amd64` | `ubuntu-latest` | `.tar.xz` | binary + `.desktop` entry | `fyne package` |
| `linux-arm64` | `ubuntu-24.04-arm` | `.tar.xz` | binary + `.desktop` entry | `fyne package` |
| `windows-amd64` | `windows-latest` | `.zip` | `Cornealius-Eyeworth.exe` | `go build` (see below) |
| `macos-arm64` | `macos-latest` | `.zip` | `Cornealius Eyeworth.app` | `fyne package` |
| `macos-amd64` | `macos-15-intel` | `.zip` | `Cornealius Eyeworth.app` | `fyne package` |

Every archive also carries `README.md`, `LICENSE.txt` and
`THIRD-PARTY-NOTICES.txt`, built per job by `scripts/gen-release-docs.sh`

## How the pipeline is wired

```
plan ── checks ─┬─ test-windows ─── build-windows ──┐
                │                                   │
                └─ test-unix ─┬──── build-linux ────┼─ vulncheck ── release
                              └──── build-macos ────┘               (tags only)
```

Each build job is a matrix over its architectures, so one job per OS but a
separate runner per target, all in parallel.

Lint first: it takes seconds, so nothing else spends a runner on a tree that does
not gofmt. govulncheck last, gating publishing without delaying the builds.

Tests gate builds, split by build constraint family - `notifier_windows.go` and
`systraypatch_win.go` compile only on Windows, their counterparts only off it. One
run per family covers both paths.

Artifacts upload on every push, so `rc` branches produce downloadable archives.
Only the release job is tag-gated.

Every Go job first rebuilds the patched Fyne fork through the `go-setup` composite
action; the `replace` in `go.mod` makes even `go vet` fail without it.

`vulncheck` lists `checks` in its needs although `checks` is already upstream: the
guard reads direct needs only, and a build skipped by a lint failure looks
identical to one skipped by a toggle. Drop that edge and lint stops blocking
releases.

## Build toggles

Targets, and whether a vulnerability finding blocks publishing, are switched in one
block in the `plan` job:

| Toggle | Default |
|--------|---------|
| `BUILD_WINDOWS_AMD64` | on |
| `BUILD_LINUX_AMD64` | on |
| `BUILD_LINUX_ARM64` | on |
| `BUILD_MACOS_ARM64` | on |
| `BUILD_MACOS_AMD64` | on |
| `VULNCHECK_BLOCKING` | on |

Turning a build off skips its job, not its tests - those are the only coverage of
the windows / non-windows split. With every architecture in one OS off, that build
job is skipped rather than handed an empty matrix, which Actions rejects.

## Testing the pipeline

Pushes to any `rc*` branch run the pipeline, as does the "Run workflow" button.
Publishing is tag-gated, so a branch run produces artifacts only - 7-day retention,
labelled with the commit SHA. It builds exactly what a tag will, so a dry run
rehearses the release rather than a subset. Narrow or remove `on.push.branches`
once releases are cut from tags alone.

## macOS

`macos-latest` is Apple Silicon and will not launch on an Intel Mac, so both
architectures ship separately. `macos-15-intel` is the last x86_64 image Actions
offers, retiring August 2027; after that Intel users build from source
(`BUILD.md`). Both bill at 10x the Linux rate on a private repo.

**The Intel archive is unverified** - there is no Intel Mac here to launch it on.

The `.app` is ad-hoc signed (`codesign -s -`) after packaging, because Apple
Silicon will not execute an unsigned arm64 binary. Not a Gatekeeper fix - users
still meet the unverified-developer dialog the README covers.

## Why Windows uses `go build` instead of `fyne package`

`rsrc_windows_amd64.syso` (`make winres`) already carries the icon and manifest;
`fyne package` adds a second one and the linker rejects the duplicate.

Winres stays the source deliberately: `winres/winres.json` declares per-monitor-v2
DPI awareness, common-controls v6 and a Win10 minimum, none of which fyne's
template has. Its filename is itself a build constraint, so a `windows-arm64`
target would ship with no icon and no manifest at all.

## The Fyne patch in CI

`go.mod` points Fyne at `../fyne-v2-watchtheme-patch`, absent from a fresh clone,
so every build rebuilds it: fetch upstream at `FYNE_TAG`, copy
`_external-patches/fyne-theme_windows.go` into place. The build fails unless the
fork differs from upstream in exactly that one file. `FYNE_TAG` is one of the four
places the version is written - `BUILD.md` lists them.

**The fork is not checksummed.** The local `replace` means Go never downloads Fyne,
so `go mod tidy` drops it from `go.sum`: this dependency sits outside supply-chain
verification, and a tag can be moved. The one-file assertion proves the tree is
unmodified, not which tree it is. A SHA in `FYNE_TAG` would close that - the
workflow already uses `git init` + `fetch` rather than `clone` so it can take one.
