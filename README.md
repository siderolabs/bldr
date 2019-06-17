# bldr

## Roadmap

- Tests
- Allow global variables
- Link using rpath or static binaries
- Dependency resolution
- Leverage labels for things like:
  - Switching the base dir install command (we default to alpine apk install for now).
    We could add a label on the bldr container `bldr.io.base.distro=[alpine,ubuntu,centos,etc.]`
  - Automatically detecting `to` in a dependency.
    We can label the container on a build with what the `finalize.to` was set to, and then automatically `COPY` from that location.
- LLB implementation
- Subpackages
  - Allow for packaging `include` and `lib` into dedicated packages
