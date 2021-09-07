# Current steps for a release

1. Update CHANGELOG.md
2. Tag
3. Push everything
4. Wait for the github workflows to complete
5. Download and verify the correct version of one of the binaries
6. Finish the draft release and publish.
7. Check gotop-builder for a successful everything build; if successful, publish.
8. Notify Nix
9. ~~Notify Homebrew~~ ~~Automated now.~~ Automation broke. Notify manually.
10. Do the Arch release.
	1. cd actions/arch-package
	2. VERSION=v4.1.2 ./run.sh

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
