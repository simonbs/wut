# Release Process for `simonbs/wut`

Follow these steps whenever you publish a new CLI version:

1. **Bump the CLI version.**
   - Update `cmd/wut/main.go` so `const version` matches the new release tag (e.g., `0.1.3`).
   - Run `go test ./...` (or the repo's agreed-upon test command) and commit the change before tagging.

2. **Create the GitHub release.**
   - Tag the commit: `git tag -a vX.Y.Z -m "Release vX.Y.Z"`.
   - Push the tag (`git push origin vX.Y.Z`) and `main` if needed.
   - Use `gh release create vX.Y.Z --title "vX.Y.Z" --notes-file /tmp/release-notes.md` (or similar `--notes`) to publish release notes describing the changes.

3. **Update the Homebrew formula.**
   - Download the release tarball (`curl -L https://github.com/simonbs/wut/archive/refs/tags/vX.Y.Z.tar.gz | shasum -a 256`) and record the SHA256.
   - Update `Formula/wut.rb` to point to the new `url` and `sha256`, then commit those changes.
   - Push the updated formula (`git push origin main`).

Optional: after publishing, confirm Homebrew builds from the new formula with `brew install --build-from-source ./Formula/wut.rb` (via a tap if required) or `brew reinstall --build-from-source wut`.
