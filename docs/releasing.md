Current steps for a release:

### gotop
1. Update Version in main.go 
2. Update CHANGELOG.md
3. Tag
4. Push everything
5. ./make.sh
6. Create github release

### Homebrew
1. Change homebrew-gotop
```
curl --output - -L https://github.com/xxxserxxx/gotop/releases/download/v3.3.2/gotop_3.3.2_linux_amd64.tgz | sha256sum
curl --output - -L https://github.com/xxxserxxx/gotop/releases/download/v3.3.2/gotop_3.3.2_darwin_amd64.tgz | sha256sum
```
