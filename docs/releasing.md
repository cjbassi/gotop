# Current steps for a release

1. Update Version in main.go 
2. Update CHANGELOG.md
3. Update docs/release.svg
4. Tag
5. Push everything
6. When the github workflows complete, finish the draft release and publish.
7. Wait for the [AUR](https://github.com/xxxserxxx/gotop-linux] project to finish building.
    1. update arch (gotop-linux) and run `aurpublish aur` and `aurpublish aur-bin`
    2. notify Nix
    3. notify Homebrew

The AUR project still needs secret credentials to aurpublish to the AUR
repository, so the final publish step is still currently manual.

Oh, what a tangled web.


## Nix 

I haven't yet figured this out, so currently just file a ticket and hope somebody on that end updates the package.

Nix adds new and interesting complexities to the release.

0. Download the gotop src package; run sha256 on it to get the hash
1. cd to the nixpkgs directory
2. Update the sha256 hash in `pkgs/tools/system/gotop/default.nix`
2. `docker run -it --rm --mount type=bind,source="\$(pwd)",target=/mnt nixos/nix sh`
3. `cd /mnt`
8. install & run vgo2nix to update deps.nix
7. `nix-build -A gotop`
8. When it fails, copy the hash and update the 


For plugin development:
```
V=$(git show -s --format=%cI HEAD | cut -b -19 |  tr -cd '[:digit:]')-$(git rev-parse HEAD | cut -b -12)
go build -ldflags "-X main.Version=$V" -o gotop ./cmd/gotop
```
