#!/usr/bin/env bash
# Builds the doc set shipped inside every release archive
set -eu

out=${1:?usage: gen-release-docs.sh <outdir>}
START_YEAR=2025
mkdir -p "$out"

cp README.md "$out/"

# --------- LICENCE ---------
year=$(date -u +%Y)
if [ "$year" = "$START_YEAR" ]; then span="$START_YEAR"; else span="$START_YEAR-$year"; fi
sed -E "s/^(Copyright \(c\) )[0-9]{4}(-[0-9]{4})?/\1$span/" LICENSE.txt > "$out/LICENSE.txt"
# A silently unmatched pattern would ship a stale year forever.
grep -q "^Copyright (c) $span " "$out/LICENSE.txt" || {
  echo "::error::LICENSE.txt copyright line did not match the expected shape"
  exit 1
}

# --------- Third-party ---------
# Platform specific
notices="$out/THIRD-PARTY-NOTICES.txt"
mods=$(go list -deps -f '{{with .Module}}{{.Path}}	{{.Version}}	{{.Dir}}{{end}}' . |
  grep . | sort -u)

self=$(go list -m)
{
  echo "Cornealius Eyeworth bundles the third-party software listed below."
  echo "Generated $(date -u +%Y-%m-%d) for $(go env GOOS)/$(go env GOARCH)."
} > "$notices"

missing=""
count=0
while IFS="	" read -r path version dir; do
  [ -z "$path" ] && continue
  [ "$path" = "$self" ] && continue
  lic=$(find "$dir" -maxdepth 1 -type f \
    \( -iname 'LICENSE*' -o -iname 'COPYING*' -o -iname 'NOTICE*' \) |
    sort | head -1)
  if [ -z "$lic" ]; then missing="$missing $path"; continue; fi
  {
    echo
    echo "=============================================================================="
    echo "$path${version:+ $version}"
    echo "=============================================================================="
    echo
    cat "$lic"
  } >> "$notices"
  count=$((count + 1))
done <<EOF
$mods
EOF

if [ -n "$missing" ]; then
  echo "::error::no licence file found for:$missing"
  exit 1
fi

echo "Release docs in $out/ - notices cover $count modules, copyright $span."
