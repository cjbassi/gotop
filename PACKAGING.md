Packaging Go for Release
========================

The gotop project in github has build rules that should compile and build a release. For development purposes, it's useful to run this (and verify the success of all the cross-compiling parts) before pushing changes to github. These are instructions on how to do this.

Dependencies
------------

All of the compiling tooling -- nearly all of which is due to cross-compiling for different platform support -- is contained in a [github actions repository](https://github.com/xxxserxxx/actions.git). You will need to check that out; this is from where you'll be building.

You will need docker or podman (I use podman, so all my examples will be podman commands).

Getting Started
---------------

The actions repo contains a README that describes how the cross-compiler works; anything that looks esoteric here (environment variables, and their legal values) is explained in that file. However, you should be able to run these commands without reading that document, to start.

Two scripts in that repo are for local use: `rebuild.sh`, and `run.sh`. The other top-level script, `entrypoint.sh` is for the container that gets built. You'll mostly be using `run.sh`.

### Step 1

- Check out the gotop repo; do *not* CD into it.
- In the repository parent's directory, start a git code server
  ```
  git daemon --port=8880 --verbose --export-all --reuseaddr --base-path=$(pwd)
  ```
  I don't fork it; I just run it in a terminal and open a different terminal for the rest.
- Check out the [github actions repo](https://github.com/xxxserxxx/actions.git), and CD into it.
- In there, run:
  ```
  ./run.sh git://localhost:8880/gotop ./cmd/gotop 'darwin/amd64 linux/amd64' refs/remotes/origin/master
  ```

It's important that the git server is running in the repo's parent directory because the build script expects a certain format to URLs to determine the project's name (among other things) -- so the project name (`gotop`) has to be in the git URL.

The first argument is that git URL, it points back to the git server you ran earlier.
The second argument is the executable path to be built. In the gotop repo, the executable is `gotop/cmd/gotop/main.go`, so that argument (relative to the project directory) is `./cmd/gotop`.

The third argument (space-separated and quoted, so it's treated as a single argument) is a list of targets to be built. gotop supports:

- darwin/amd64/1
- darwin/arm64/1
- linux/amd64
- linux/386
- linux/arm64
- linux/arm7
- linux/arm6
- linux/arm5
- windows/amd64/1
- windows/386/1
- freebsd/amd64/1

The structure parts are parsed into: `$GOOS/$GOARCH/$CGO`; you can mix and match how you want, and add different GOOS and GOARCH combinations, and change the CGO value -- YMMV. It might compile. It may even run. But that list is what gotop officially supports. The only optional part is the `CGO` parameter; it defaults to 0.

The last argument is the git reference to a branch. If you're building master, just use the one in the example; otherwise, figure it out for yourself because I can't explain it for you -- I'm really a Mercurial guy who only uses git and github when I'm forced to, and gotop's original author started it in github. ¯\_(ツ)_/¯

The first time you build, a bunch of containers will be downloaded and built. The script tries to be smart about reusing build containers, so subsequent runs should run faster.

Artifacts will be in `work/gotop/gotop/.release`.

Erratta
-------

The build process is both limited in some ways, and has more features in others. It grew up around gotop, and so inherits some assumptions from that project; I tried to generalize it, but it still has some gotop biases. If you read the scripts (they're all shell scripts) there are some command-line arguments that let you do some things, like rebuild the compile container (e.g., if you change `entrypoint.sh`); I leave that as a voyage of discovery for the terminally curious.

I *happily* accept pull requests to improve the scripts. Bug fixes, more capabilities, removing the gotop biases (making more generic for use with other Go projects), spelling corrections, whatever.

If you want to see an example of how this is used in the gotop project, check out the [the workflows](https://github.com/xxxserxxx/gotop/tree/master/.github/workflows).

Good luck, and if you have any questions, pop into [#gotop:matrix.org](https://app.element.io/#/room/#gotop:matrix.org) with your favorite Matrix client and give me a shout (by name). It may take me a while to respond (maybe even a couple days), but I *will* respond.
