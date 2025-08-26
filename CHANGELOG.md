## [bldr 0.5.2](https://github.com/siderolabs/bldr/releases/tag/v0.5.2) (2025-08-26)

Welcome to the v0.5.2 release of bldr!



Please try out the release binaries and report any issues at
https://github.com/siderolabs/bldr/issues.

### Contributors

* Andrey Smirnov
* Dmitrii Sharshakov
* Mateusz Urbanek

### Changes
<details><summary>3 commits</summary>
<p>

* [`73b4eaa`](https://github.com/siderolabs/bldr/commit/73b4eaa7d798f1afe460f9114b9284bde6e0efd9) feat: implement a dump command
* [`42e5c02`](https://github.com/siderolabs/bldr/commit/42e5c0290f6aca41a6ab3a8f52b970f360c81e2d) feat: add support for automatic templating
* [`c517c50`](https://github.com/siderolabs/bldr/commit/c517c50c23888093d94de462010e31158a52c830) release(v0.5.1): prepare release
</p>
</details>

### Changes from siderolabs/gen
<details><summary>1 commit</summary>
<p>

* [`044d921`](https://github.com/siderolabs/gen/commit/044d921685bbd8b603a64175ea63b07efe9a64a7) feat: add xslices.Deduplicate
</p>
</details>

### Dependency Changes

* **github.com/emicklei/dot**          v1.8.0 -> v1.9.0
* **github.com/google/go-github/v74**  v74.0.0 **_new_**
* **github.com/siderolabs/gen**        v0.8.4 -> v0.8.5
* **golang.org/x/sync**                v0.15.0 -> v0.16.0

Previous release can be found at [v0.5.1](https://github.com/siderolabs/bldr/releases/tag/v0.5.1)

## [bldr 0.5.1](https://github.com/siderolabs/bldr/releases/tag/v0.5.1) (2025-07-15)

Welcome to the v0.5.1 release of bldr!



Please try out the release binaries and report any issues at
https://github.com/siderolabs/bldr/issues.

### Contributors

* Dmitrii Sharshakov

### Changes
<details><summary>1 commit</summary>
<p>

* [`3ea108b`](https://github.com/siderolabs/bldr/commit/3ea108b0dce619891052792c4dcb5b179109ca42) feat: add `name` parameter to override names in SBOM
</p>
</details>

### Dependency Changes

This release has no dependency changes

Previous release can be found at [v0.5.0](https://github.com/siderolabs/bldr/releases/tag/v0.5.0)


## [bldr 0.5.0](https://github.com/siderolabs/bldr/releases/tag/v0.5.0) (2025-07-09)

Welcome to the v0.5.0 release of bldr!



Please try out the release binaries and report any issues at
https://github.com/siderolabs/bldr/issues.

### Contributors

* Dmitrii Sharshakov
* Artem Chernyshev
* Andrey Smirnov
* Dmitriy Matrenichev
* Utku Ozdemir

### Changes
<details><summary>4 commits</summary>
<p>

* [`dafc274`](https://github.com/siderolabs/bldr/commit/dafc274b9454472233f7e4788195d5fc47a6c57f) fix: make SBOM generation work in a multi-step build
* [`38d5691`](https://github.com/siderolabs/bldr/commit/38d5691bb40874088e734aa4c22862166017ca9f) fix: sbom: do not append arch
* [`e45241d`](https://github.com/siderolabs/bldr/commit/e45241de914f4be03ee8311005e1f2d380b1fa2f) feat: add SBOM support
* [`d09e69c`](https://github.com/siderolabs/bldr/commit/d09e69caf574060feb26ae8ead006b9161da8e5c) feat: update dependencies
</p>
</details>

### Changes from siderolabs/gen
<details><summary>4 commits</summary>
<p>

* [`dcb2b74`](https://github.com/siderolabs/gen/commit/dcb2b7417879f230a569ce834dad5c89bd09d6bf) feat: add `panicsafe` package
* [`b36ee43`](https://github.com/siderolabs/gen/commit/b36ee43f667a7a56b340a3e769868ff2a609bb5b) feat: make `xyaml.CheckUnknownKeys` public
* [`3e319e7`](https://github.com/siderolabs/gen/commit/3e319e7e52c5a74d1730be8e47952b3d16d91148) feat: implement `xyaml.UnmarshalStrict`
* [`7c0324f`](https://github.com/siderolabs/gen/commit/7c0324fee9a7cfbdd117f43702fa273689f0db97) chore: future-proof HashTrieMap
</p>
</details>

### Dependency Changes

* **github.com/anchore/syft**               v1.27.1 **_new_**
* **github.com/emicklei/dot**               v1.6.4 -> v1.8.0
* **github.com/moby/buildkit**              v0.19.0 -> v0.23.0
* **github.com/opencontainers/image-spec**  v1.1.0 -> v1.1.1
* **github.com/siderolabs/gen**             v0.8.0 -> v0.8.4
* **github.com/spf13/cobra**                v1.8.1 -> v1.9.1
* **golang.org/x/oauth2**                   v0.26.0 -> v0.30.0
* **golang.org/x/sync**                     v0.11.0 -> v0.15.0

Previous release can be found at [v0.4.1](https://github.com/siderolabs/bldr/releases/tag/v0.4.1)

## [bldr 0.4.1](https://github.com/siderolabs/bldr/releases/tag/v0.4.1) (2025-02-06)

Welcome to the v0.4.1 release of bldr!



Please try out the release binaries and report any issues at
https://github.com/siderolabs/bldr/issues.

### Contributors

* Andrey Smirnov
* Dmitriy Matrenichev
* Noel Georgi

### Changes
<details><summary>9 commits</summary>
<p>

* [`b2a223d`](https://github.com/siderolabs/bldr/commit/b2a223d14d48b88e39ead81772953efc2f4dbdb5) fix: exclude test script output from the result
* [`76a2c8f`](https://github.com/siderolabs/bldr/commit/76a2c8f404b546e5d1d1088c3a0e2406d371ff6d) fix: use single build directory
* [`62ed470`](https://github.com/siderolabs/bldr/commit/62ed470038fb769502600ffb5a3630bbb20c58f5) release(v0.4.0): prepare release
* [`4a79aeb`](https://github.com/siderolabs/bldr/commit/4a79aebb10957f23c7fc82378b4498686bb6fe73) feat: make variant in `pkg.yaml` explicit
* [`d978bcc`](https://github.com/siderolabs/bldr/commit/d978bcc7571839ff65d8a139be600c57135cdbf6) fix: drop TOOLCHAIN from std vars
* [`66b078d`](https://github.com/siderolabs/bldr/commit/66b078d12e99b92c7e9b3a9ceee71936ad52bd4e) fix: drop alpine for `md5sum` step
* [`24a37b0`](https://github.com/siderolabs/bldr/commit/24a37b0f070182ce8706e852f22f54f799761f33) feat: implement support for network mode
* [`40ad501`](https://github.com/siderolabs/bldr/commit/40ad5014b4fc1a187d2c83e1ebcd18f40890a24b) chore: rekres to simplify `.kres.yaml` defaults
* [`73fad60`](https://github.com/siderolabs/bldr/commit/73fad60813195cef4458685335e5e907726bd21e) chore: bump deps
</p>
</details>

### Changes since v0.4.0
<details><summary>2 commits</summary>
<p>

* [`b2a223d`](https://github.com/siderolabs/bldr/commit/b2a223d14d48b88e39ead81772953efc2f4dbdb5) fix: exclude test script output from the result
* [`76a2c8f`](https://github.com/siderolabs/bldr/commit/76a2c8f404b546e5d1d1088c3a0e2406d371ff6d) fix: use single build directory
</p>
</details>

### Changes from siderolabs/gen
<details><summary>4 commits</summary>
<p>

* [`5ae3afe`](https://github.com/siderolabs/gen/commit/5ae3afee65490ca9f4bd32ea41803ab3a17cad7e) chore: update hashtriemap implementation from the latest upstream
* [`e847d2a`](https://github.com/siderolabs/gen/commit/e847d2ace9ede4a17283426dfbc8229121f2909b) chore: add more utilities to xiter
* [`f3c5a2b`](https://github.com/siderolabs/gen/commit/f3c5a2b5aba74e4935d073a0135c4904ef3bbfef) chore: add `Empty` and `Empty2` iterators
* [`c53b90b`](https://github.com/siderolabs/gen/commit/c53b90b4a418b8629d938af06900249ce5acd9e6) chore: add packages xiter/xstrings/xbytes
</p>
</details>

### Dependency Changes

* **github.com/Masterminds/sprig/v3**  v3.2.3 -> v3.3.0
* **github.com/containerd/platforms**  v0.2.1 -> v1.0.0-rc.1
* **github.com/emicklei/dot**          v1.6.2 -> v1.6.4
* **github.com/google/go-github/v67**  v67.0.0 **_new_**
* **github.com/moby/buildkit**         v0.15.1 -> v0.19.0
* **github.com/otiai10/copy**          v1.14.0 -> v1.14.1
* **github.com/siderolabs/gen**        v0.5.0 -> v0.8.0
* **github.com/stretchr/testify**      v1.9.0 -> v1.10.0
* **golang.org/x/oauth2**              v0.22.0 -> v0.26.0
* **golang.org/x/sync**                v0.8.0 -> v0.11.0

Previous release can be found at [v0.3.2](https://github.com/siderolabs/bldr/releases/tag/v0.3.2)

## [bldr 0.4.0](https://github.com/siderolabs/bldr/releases/tag/v0.4.0) (2025-01-31)

Welcome to the v0.4.0 release of bldr!



Please try out the release binaries and report any issues at
https://github.com/siderolabs/bldr/issues.

### Contributors

* Dmitriy Matrenichev
* Andrey Smirnov
* Noel Georgi

### Changes
<details><summary>6 commits</summary>
<p>

* [`4a79aeb`](https://github.com/siderolabs/bldr/commit/4a79aebb10957f23c7fc82378b4498686bb6fe73) feat: make variant in `pkg.yaml` explicit
* [`d978bcc`](https://github.com/siderolabs/bldr/commit/d978bcc7571839ff65d8a139be600c57135cdbf6) fix: drop TOOLCHAIN from std vars
* [`66b078d`](https://github.com/siderolabs/bldr/commit/66b078d12e99b92c7e9b3a9ceee71936ad52bd4e) fix: drop alpine for `md5sum` step
* [`24a37b0`](https://github.com/siderolabs/bldr/commit/24a37b0f070182ce8706e852f22f54f799761f33) feat: implement support for network mode
* [`40ad501`](https://github.com/siderolabs/bldr/commit/40ad5014b4fc1a187d2c83e1ebcd18f40890a24b) chore: rekres to simplify `.kres.yaml` defaults
* [`73fad60`](https://github.com/siderolabs/bldr/commit/73fad60813195cef4458685335e5e907726bd21e) chore: bump deps
</p>
</details>

### Changes from siderolabs/gen
<details><summary>4 commits</summary>
<p>

* [`5ae3afe`](https://github.com/siderolabs/gen/commit/5ae3afee65490ca9f4bd32ea41803ab3a17cad7e) chore: update hashtriemap implementation from the latest upstream
* [`e847d2a`](https://github.com/siderolabs/gen/commit/e847d2ace9ede4a17283426dfbc8229121f2909b) chore: add more utilities to xiter
* [`f3c5a2b`](https://github.com/siderolabs/gen/commit/f3c5a2b5aba74e4935d073a0135c4904ef3bbfef) chore: add `Empty` and `Empty2` iterators
* [`c53b90b`](https://github.com/siderolabs/gen/commit/c53b90b4a418b8629d938af06900249ce5acd9e6) chore: add packages xiter/xstrings/xbytes
</p>
</details>

### Dependency Changes

* **github.com/Masterminds/sprig/v3**  v3.2.3 -> v3.3.0
* **github.com/containerd/platforms**  v0.2.1 -> v1.0.0-rc.1
* **github.com/emicklei/dot**          v1.6.2 -> v1.6.4
* **github.com/google/go-github/v67**  v67.0.0 **_new_**
* **github.com/moby/buildkit**         v0.15.1 -> v0.19.0
* **github.com/otiai10/copy**          v1.14.0 -> v1.14.1
* **github.com/siderolabs/gen**        v0.5.0 -> v0.8.0
* **github.com/stretchr/testify**      v1.9.0 -> v1.10.0
* **golang.org/x/oauth2**              v0.22.0 -> v0.25.0
* **golang.org/x/sync**                v0.8.0 -> v0.10.0

Previous release can be found at [v0.3.2](https://github.com/siderolabs/bldr/releases/tag/v0.3.2)

## [bldr 0.2.2](https://github.com/siderolabs/bldr/releases/tag/v0.2.2) (2023-09-15)

Welcome to the v0.2.2 release of bldr!



Please try out the release binaries and report any issues at
https://github.com/siderolabs/bldr/issues.

### Contributors

* Noel Georgi
* Andrey Smirnov

### Changes
<details><summary>3 commits</summary>
<p>

* [`20de986`](https://github.com/siderolabs/bldr/commit/20de986bf1d15d50bd05cf62da99266af0b590a5) feat: support --cache-from/--cache-imports
* [`38bee98`](https://github.com/siderolabs/bldr/commit/38bee98982471791133abaaa89dbffe225adeb3a) chore: rekres to support gh actions
* [`d7f236f`](https://github.com/siderolabs/bldr/commit/d7f236f3fb0d91382a328b54ebe791b59c9e1f45) chore: add a no-op github workflow
</p>
</details>

### Dependency Changes

* **github.com/containerd/containerd**  v1.7.3 -> v1.7.6
* **github.com/moby/buildkit**          v0.12.1 -> v0.12.2
* **golang.org/x/oauth2**               v0.11.0 -> v0.12.0

Previous release can be found at [v0.2.1](https://github.com/siderolabs/bldr/releases/tag/v0.2.1)




## [bldr 0.2.0](https://github.com/siderolabs/bldr/releases/tag/v0.2.0) (2023-05-19)

Welcome to the v0.2.0 release of bldr!



Please try out the release binaries and report any issues at
https://github.com/siderolabs/bldr/issues.

### Contributors

* Andrey Smirnov
* Noel Georgi
* Andrey Smirnov
* Andrew Rynhard
* Alexey Palazhchenko
* Alexey Palazhchenko
* Artem Chernyshev
* Dmitriy Matrenichev
* Andrey Smirnov
* Brad Beam
* Niklas Wik

### Changes
<details><summary>81 commits</summary>
<p>

* [`97650b2`](https://github.com/siderolabs/bldr/commit/97650b2f1008755eb96718b02506da6ec9a9d5da) feat: implement --no-cache option
* [`bc2438e`](https://github.com/siderolabs/bldr/commit/bc2438e6c434ab5ab1d4b1c2fe85267807a062fd) chore: bump deps
* [`f8905e8`](https://github.com/siderolabs/bldr/commit/f8905e8891415c42f43a37b7d31e3198ca41148f) chore: bump deps
* [`4c077ad`](https://github.com/siderolabs/bldr/commit/4c077ad50b9f2bb86eba31ab0e38c27d6e98776d) chore: rekres + bump deps
* [`3d32358`](https://github.com/siderolabs/bldr/commit/3d323585edc20f7da6ce97ebf19f2e021dc799e1) chore: bump deps
* [`a7421af`](https://github.com/siderolabs/bldr/commit/a7421af3d6f3f72dec684151e89b29990e62ac5b) chore: bump deps
* [`4f20353`](https://github.com/siderolabs/bldr/commit/4f20353d5e3ae974b3fe2a9cbd0cd227bc751a36) chore: bump deps
* [`ed832fe`](https://github.com/siderolabs/bldr/commit/ed832fee42145eac7f50a10509385af113dceda6) chore: bump deps
* [`f2f1d18`](https://github.com/siderolabs/bldr/commit/f2f1d18f5973fa26aaeaba16852164530976c408) chore: bump dependencies
* [`77bd646`](https://github.com/siderolabs/bldr/commit/77bd64695d8fb7c299770ebc885cbfcc0b8aa34e) chore: rekres
* [`0af77c6`](https://github.com/siderolabs/bldr/commit/0af77c6a71ebd007dabe56b44852463cf26a29a0) chore: bump go dependencies
* [`a46c868`](https://github.com/siderolabs/bldr/commit/a46c86846953ba25f67c8ebf1d5261c6bb55e4dc) chore: bump dependencies
* [`804a03b`](https://github.com/siderolabs/bldr/commit/804a03ba8d2ab379f0f1a3ecca708b9e45e11b17) chore: bump dependencies
* [`1218936`](https://github.com/siderolabs/bldr/commit/12189362f4fccac5205b9f527dbbc402d3269dce) chore: bump deps
* [`afa3012`](https://github.com/siderolabs/bldr/commit/afa301220eaebdea9ef87ef22d5e2e389f4e8c95) chore: disable drone slack pipeline for renovate
* [`f2690f0`](https://github.com/siderolabs/bldr/commit/f2690f0a188e78d91d6c1c9d7619ac13e99c7367) chore: disable drone for renovate/dependabot
* [`cb072a4`](https://github.com/siderolabs/bldr/commit/cb072a42d0bd61fadf21158c7c5d9791c11b4934) chore: bump deps
* [`42c0ffb`](https://github.com/siderolabs/bldr/commit/42c0ffb7ea8048aa353a71aaa2edae994041f150) fix: multi arch build
* [`5613733`](https://github.com/siderolabs/bldr/commit/5613733b0afc35b815b8cb6413ab4c9ad0a33b10) chore: update renovate config
* [`9b2810f`](https://github.com/siderolabs/bldr/commit/9b2810f5628b91e5a4b467f6c9b3dc0894bfab13) chore: kresify, rename
* [`b979fbd`](https://github.com/siderolabs/bldr/commit/b979fbd2ee8943189ba7d1324d860c24ebb438a2) feat: implement additional variables and `bldr eval`
* [`bf9bada`](https://github.com/siderolabs/bldr/commit/bf9bada3887973c2d3747ca0dc260ab70060aa2f) fix: remove "promoted dependency" feature
* [`a714657`](https://github.com/siderolabs/bldr/commit/a7146576d131c2cd02950dbbc0a5196e8c43b74d) fix: skip merging if no dependencies found
* [`cedc88a`](https://github.com/siderolabs/bldr/commit/cedc88afad45bd9fa9806898c70569b60952fc68) feat: support build behind proxy
* [`81055c0`](https://github.com/siderolabs/bldr/commit/81055c0ec29ef7e59c38888d5f8d4387199af4bf) fix: make sure cache persists local context changes
* [`03757bb`](https://github.com/siderolabs/bldr/commit/03757bbefd06252193e1aacf56b3e05fa39c112a) feat: use new LLB merge operation, allow `from:` in deps
* [`febf1d0`](https://github.com/siderolabs/bldr/commit/febf1d058184ee5c40694b25515685847a43b1da) chore: bump dependencies
* [`9d49478`](https://github.com/siderolabs/bldr/commit/9d49478dc112cd1badedb0d02e190dc004b6fb62) chore: rename org
* [`376fe2b`](https://github.com/siderolabs/bldr/commit/376fe2b51afaf7f5d6b03a22bf550fce751319c9) feat: implement `SOURCE_DATE_EPOCH` build argument
* [`7a0ad18`](https://github.com/siderolabs/bldr/commit/7a0ad18e3627300e29ce43088a64735b964a2d3e) feat: support cache mounts
* [`86cfe12`](https://github.com/siderolabs/bldr/commit/86cfe1251e0ed61758ed2483087b4842df6918c4) fix: respect HTTP proxy settings
* [`e2c007a`](https://github.com/siderolabs/bldr/commit/e2c007a86bcc7c63b60e336b5cbe13f47529276a) feat: optimize build time by removing unnecessary I/O
* [`384f28d`](https://github.com/siderolabs/bldr/commit/384f28d53e6b46cb28f836c83d7d3d200959b25b) chore: add debug `--json` flag to `llb`
* [`d4df177`](https://github.com/siderolabs/bldr/commit/d4df177f8c458ab672f2817b1352603cff7fd4ff) feat: bump Alpine image for the `alpine` image
* [`6ffa6b2`](https://github.com/siderolabs/bldr/commit/6ffa6b27d6e90ca069d2cd107927d343741ba082) feat: add destination name to error messages
* [`3198175`](https://github.com/siderolabs/bldr/commit/3198175d11e21abbc1982ef4efeed45acd817f20) chore: bump dependencies
* [`8c545bd`](https://github.com/siderolabs/bldr/commit/8c545bda7798f92c1458391eec5fd0e627980dbf) fix: detect updates for containerd
* [`a8d27f7`](https://github.com/siderolabs/bldr/commit/a8d27f75d69a2b7011e2f7bd82d8df5708b2e14d) fix: detect updates for runc and similar pkgs
* [`aa62d4a`](https://github.com/siderolabs/bldr/commit/aa62d4a95b363f7cd893287d94dd8859e266f0ee) feat: add checksum validation
* [`07cd6ea`](https://github.com/siderolabs/bldr/commit/07cd6eae82156960969b383e8262767eb6318a87) feat: implement `bldr update --dry` command
* [`533e360`](https://github.com/siderolabs/bldr/commit/533e360b70b7355fb6ce39409d182b05ae62f1df) feat: add support for image labels
* [`f27a804`](https://github.com/siderolabs/bldr/commit/f27a804661d14d11f9080ed1a9d4a5fe26923279) fix: support multi-platform key passed in via options
* [`f3b2dcf`](https://github.com/siderolabs/bldr/commit/f3b2dcfb1824cef1b64380f16aae7a0758773d8b) fix: update properly reference to the image in export map
* [`f71d92d`](https://github.com/siderolabs/bldr/commit/f71d92dd7bdaac547dc73b11ad1abea2be265ece) chore: fix master build by adding `SSH_KEY`
* [`e42dfc4`](https://github.com/siderolabs/bldr/commit/e42dfc41f8c6257f95c7441dead74c0f0948f28a) feat: build bldr for linux/amd64 and linux/arm64
* [`39b6665`](https://github.com/siderolabs/bldr/commit/39b6665ea6be1ad25c9b3ea781a7422f24d33f6f) docs: use correct term instead of 'shebang'
* [`47a36b3`](https://github.com/siderolabs/bldr/commit/47a36b363db4ed4da426d9c3a6168c1772d53b35) fix: linter shouldn't complain on `finalize` on `steps` missing
* [`4f43f7e`](https://github.com/siderolabs/bldr/commit/4f43f7ef3b3f091de687feb2a2d7b65ad5bcca3c) chore: use git instead of gitmeta
* [`7745285`](https://github.com/siderolabs/bldr/commit/7745285a1fc6f0b652eb772eb1a5cb05505dc846) test: add integration test for variables and Sprig template funcs
* [`ba41ffb`](https://github.com/siderolabs/bldr/commit/ba41ffba28240cb282d61313157a8fb5539540ba) chore: fancy slack webhook message
* [`b40ae96`](https://github.com/siderolabs/bldr/commit/b40ae96d78d16a5b4a94131d5a731a6ec80175c6) feat: implement integration tests for bldr
* [`fb7f613`](https://github.com/siderolabs/bldr/commit/fb7f613c141e698e0058f226d68cd487f5274878) feat: add `bldr validate` command
* [`31a7421`](https://github.com/siderolabs/bldr/commit/31a74218aa7f0424e4eff4f8e31c87348d60151a) feat: stop on error in any package, more validation
* [`b10aef0`](https://github.com/siderolabs/bldr/commit/b10aef032c74e4d3a4cb378a0dcc530f04e3e27c) chore: move 'upgrade' under `pkg/`
* [`8513435`](https://github.com/siderolabs/bldr/commit/851343553603812ff4dd1e46f25030a5364b8a56) docs: add documentation for the `bldr`
* [`94637ff`](https://github.com/siderolabs/bldr/commit/94637ff23327495e4e7c9ed9092d9adf4e2bdc68) fix: add 'test' step to the build
* [`d8b00d4`](https://github.com/siderolabs/bldr/commit/d8b00d40620be87434049da7768c71a3c699dd45) chore: remove ca-certificates from base image
* [`1289eba`](https://github.com/siderolabs/bldr/commit/1289ebad6a228be0300c524748adfecc39ae34db) feat: implement `Pkgfile` and support variables for templates
* [`2b9366f`](https://github.com/siderolabs/bldr/commit/2b9366f80af0cb26cf99cf1bae0bb4c83194abb3) feat: add support for runtime dependencies
* [`ae7df38`](https://github.com/siderolabs/bldr/commit/ae7df38c9cef96424b1b110d5784261f6b8f642f) feat: add Alpine packages to the dependency graph
* [`067a1ea`](https://github.com/siderolabs/bldr/commit/067a1ea1db4b5c4657c3a84b8d3c836081c0fd68) feat: allow 'graph' to be generated without 'target'
* [`cdc4af3`](https://github.com/siderolabs/bldr/commit/cdc4af324b34ee0f9987a357022bfa7c045be11c) feat: graph external dependencies (images) in addition to stages
* [`3576a53`](https://github.com/siderolabs/bldr/commit/3576a53b846cbcfdabadd02ce5b5402242d62db5) chore: add golangci-lint and fix linter issues
* [`285fcfe`](https://github.com/siderolabs/bldr/commit/285fcfe36f5295f59fefecdf4ebb6239dab4dada) chore: remove 'bldr convert' command
* [`2b3b543`](https://github.com/siderolabs/bldr/commit/2b3b5431ec2255a46ba54c4b5b545f18ad896d11) chore: bump docker version to enable image push
* [`8292442`](https://github.com/siderolabs/bldr/commit/8292442a1ad1bc8eb3222c29fe42343148d0a479) chore: make bldr push its images on 'master' build
* [`ea80f4e`](https://github.com/siderolabs/bldr/commit/ea80f4eca626478c0ed850cc523831d100c0af6f) feat: new version of bldr based on buildkit LLB
* [`e37bd48`](https://github.com/siderolabs/bldr/commit/e37bd48d007c89cf346c0f2f82921725fa32a458) chore: use Alpine 3.10 as base image for the builds
* [`2be6335`](https://github.com/siderolabs/bldr/commit/2be63354c638d48a76fe4b1ab88018f25fac8d9f) chore: default organization/registry to the ones bldr was built with
* [`62ea4a5`](https://github.com/siderolabs/bldr/commit/62ea4a53e8557f9ecba1642bddc0299ad97986e3) chore: fix race in verify method
* [`0272817`](https://github.com/siderolabs/bldr/commit/0272817d5193dcdef843b9cc6e16fb0befb9ba22) chore: update channel name
* [`78d400c`](https://github.com/siderolabs/bldr/commit/78d400cc7554e4885728c80376297a7433d2d086) chore: bump docker in Drone (#11)
* [`946e61b`](https://github.com/siderolabs/bldr/commit/946e61b2e0dd517571462af6861e4ab1ef3ec803) fix: fix typo in path (#8)
* [`5d96fb8`](https://github.com/siderolabs/bldr/commit/5d96fb877aad63708c630029c46bd3e40ded23a6) fix: set the build path last (#7)
* [`eb82bb7`](https://github.com/siderolabs/bldr/commit/eb82bb7dfe854ad9468492ffe3dea41d0a4292fe) fix: add the cache flags (#6)
* [`fe3cbd3`](https://github.com/siderolabs/bldr/commit/fe3cbd3a9239cdccb16f76d8cc92992afd3ed482) chore: set docker mtu to 1440 (#5)
* [`f5305ac`](https://github.com/siderolabs/bldr/commit/f5305acba61b24949e049c18e643f17d5336b526) feat: move cache flags to CLI options (#4)
* [`95bb0e7`](https://github.com/siderolabs/bldr/commit/95bb0e784149a0853657b0975c787ce65f0a0c6f) feat: initial implementation (#3)
* [`988e857`](https://github.com/siderolabs/bldr/commit/988e8570e4955cc7454dfb56643795e374926f6c) Merge pull request  [#1](https://github.com/siderolabs/bldr/pull/1) from andrewrynhard/initial-implementation
* [`5c25018`](https://github.com/siderolabs/bldr/commit/5c250181e8d3699a4e30716bb8112cb246853b22) feat: initial implementation
* [`7c244c6`](https://github.com/siderolabs/bldr/commit/7c244c6fbcb04fdd9d3aa1daa1f77cb01f7af93d) Initial commit
</p>
</details>

### Changes since v0.2.0-alpha.12
<details><summary>4 commits</summary>
<p>

* [`97650b2`](https://github.com/siderolabs/bldr/commit/97650b2f1008755eb96718b02506da6ec9a9d5da) feat: implement --no-cache option
* [`bc2438e`](https://github.com/siderolabs/bldr/commit/bc2438e6c434ab5ab1d4b1c2fe85267807a062fd) chore: bump deps
* [`f8905e8`](https://github.com/siderolabs/bldr/commit/f8905e8891415c42f43a37b7d31e3198ca41148f) chore: bump deps
* [`4c077ad`](https://github.com/siderolabs/bldr/commit/4c077ad50b9f2bb86eba31ab0e38c27d6e98776d) chore: rekres + bump deps
</p>
</details>

### Dependency Changes

This release has no dependency changes

