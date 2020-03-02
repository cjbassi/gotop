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
