# Contributing

## Getting Started

We'd love to accept your contributions to this project! If you are a first time contributor, please review our [Contributing Guidelines](https://go-vela.github.io/docs/community/contributing_guidelines/) before proceeding.

### Prerequisites

- [Review the commit guide we follow](https://chris.beams.io/posts/git-commit/#seven-rules) - ensure your commits follow our standards
- Review our [style guide](https://go-vela.github.io/docs/community/contributing_guidelines/#style-guide) to ensure your code is clean and consistent.
- [Docker](https://docs.docker.com/install/) - building block for local development
- [Docker Compose](https://docs.docker.com/compose/install/) - start up local development
- [Make](https://www.gnu.org/software/make/) - start up local development

### Setup

- [Fork](/fork) this repository

- Clone this repository to your workstation:

```bash
# Clone the project
git clone git@github.com:go-vela/vela-k6.git $HOME/go-vela/vela-k6
```

- Navigate to the repository code:

```bash
# Change into the project directory
cd $HOME/go-vela/vela-k6
```

- Point the original code at your fork:

```bash
# Add a remote branch pointing to your fork
git remote add fork https://github.com/your_fork/vela-k6
```

### Running Locally

- Navigate to the repository code:

```bash
# Change into the project directory
cd $HOME/go-vela/vela-k6
```

- Build the repository code:

```bash
# Build the code with `make`
make build
```

- Run the repository code:

```bash
# Run the executable
./vela-k6
```

### Development

- Navigate to the repository code:

```bash
# Change into the project directory
cd $HOME/go-vela/vela-k6
```

- Write your code and tests to implement the changes you desire.

- Run the repository code (ensures your changes perform as you desire):
  - Ensure the appropriate environment variables are set (`PARAMETER_SCRIPT_PATH`, optionally `PARAMETER_OUTPUT_PATH`, `PARAMETER_PROJEKTOR_COMPAT_MODE`, and `PARAMETER_FAIL_ON_THRESHOLD_BREACH`)

```bash
# execute the `start` target with `make`
make start
```

- Test the repository code (ensures your changes don't break existing functionality):

```bash
# execute the `test` target with `make`
make test
```

- Check the repository code (ensures your code meets the project standards):

```bash
# execute the `check` target with `make`
make check
```

- Push to your fork:

```bash
# push your code up to your fork
git push fork main
```

- Make sure to follow our [PR process](https://go-vela.github.io/docs/community/contributing_guidelines/#development-workflow) when opening a pull request

Thank you for your contribution!
