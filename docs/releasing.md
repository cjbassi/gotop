Current steps for a release:

### gotop
1. Update Version in main.go 
2. Update CHANGELOG.md
3. Tag
4. Push everything
5. When the github workflows complete, finish the draft release and publish.
6. After the [Homebrew](https://github.com/xxxserxxx/homebrew-gotop) and [AUR](https://github.com/xxxserxxx/gotop-linux] projects are done, check out gotop-linux and run `aurpublish aur` and `aurpublish aur-bin`


Homebrew is automatically updated.  The AUR project still needs secret
credentials to aurpublish to the AUR repository, so the final publish step is
still currently manual.

Oh, what a tangled web.


Nix adds new and interesting complexities to the release.

1. cd to the nixpkgs directory
2. docker run -it --rm --mount type=bind,source="\$(pwd)",target=/mnt nixos/nix sh
3. cd /mnt
4. nix-prefetch-url --unpack https://github.com/xxxserxxx/gotop/archive/v3.3.2.tar.gz
5. Copy the sha256
6. Update the version and hash in nixpkgs/pkgs/tools/system/gotop/default.nix
8. In docker, install & run vgo2nix to update deps.nix
7. nix-build -A gotop


For plugin development:
```
V=$(git show -s --format=%cI HEAD | cut -b -19 |  tr -cd '[:digit:]')-$(git rev-parse HEAD | cut -b -12)
go build -ldflags "-X main.Version=$V" -o gotop ./cmd/gotop
```
