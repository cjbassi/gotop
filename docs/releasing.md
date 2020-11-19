# Current steps for a release

1. Update Version in main.go 
2. Update CHANGELOG.md
3. Tag
4. Push everything
5. Wait for the github workflows to complete
6. Download and verify the correct version of one of the binaries
7. Finish the draft release and publish.
8. Check gotop-builder for a successful everything build; if successful, publish.
10. Wait for the [AUR](https://github.com/xxxserxxx/gotop-linux) project to finish building.
    1. update arch (gotop-linux) and run `aurpublish gotop` and `aurpublish gotop-bin`
    2. Test install `gotop` and `gotop-bin` with running & version check
11. Notify Nix
12. ~~Notify Homebrew~~ Automated now.

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
8. When it fails, ...
