# Current steps for a release

1. Update CHANGELOG.md
2. Tag
3. Push everything
4. Wait for the github workflows to complete
5. Download and verify the correct version of one of the binaries
6. Finish the draft release and publish.
7. Check gotop-builder for a successful everything build; if successful, publish.
8. ~~Notify Nix~~ -- this seems to have been automated by the Nix folks?
9. ~~Notify Homebrew~~ ~~Automated now.~~ Automation broke. Notify manually.
10. Do the Arch release.
	1. cd actions/arch-package
	2. VERSION=v4.1.2 ./run.sh
