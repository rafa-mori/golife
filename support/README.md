# 🛠️ Build, Install and Test Flow for Go Modules

This directory contains reusable support scripts for automating build, installation, compression, validation, and publishing of Go modules. The flow is fully dynamic and can be copied to any project—just set the OWNER and app name.

## 📦 Script Structure

- **config.sh**: Centralizes global variables (OWNER, APP_NAME, directories, version, etc).
- **platform.sh**: Detects OS, architecture, and builds release URLs.
- **utils.sh**: Utility functions (colored logs, temp directory handling, shell helpers).
- **install_funcs.sh**: Functions for installation, PATH management, compression, download and binary installation.
- **validate.sh**: Dependency and Go version validation.
- **info.sh**: Shows banners, project info, and install summary.
- **install.sh**: Main install script, orchestrates the flow and loads the libs.
- **build.sh**: Build script, cross-compiling, compression and packaging.

## 🚀 How to Use

1. **Set OWNER and APP_NAME** in `config.sh`.
2. **Build:**

   ```sh
   ./support/build.sh
   ```

3. **Install:**

   ```sh
   ./support/install.sh
   ```

4. **Validate dependencies:**

   ```sh
   ./support/validate.sh
   ```

## 🧩 How to Reuse in Other Projects

- Copy the `support/` folder to your new project.
- Adjust `OWNER`, `APP_NAME` and variables in `config.sh`.
- Adapt banners and messages if you wish.
- Done! The flow covers build, compression, install, PATH management and validation.

## 💡 Tips

- All logs are colored and centralized via the `log` function.
- The flow detects root/non-root and installs in the correct location.
- PATH is managed automatically and safely.
- Build is cross-platform and easily extensible.
- Scripts are Shellcheck and POSIX compatible.

## 📚 Usage Examples

See other ecosystem modules for usage and customization examples:

- [goforge](https://github.com/rafa-mori/goforge)
- [gdbase](https://github.com/rafa-mori/gdbase)
- [logz](https://github.com/rafa-mori/logz)

---

**Maintainer:** Rafael Mori

Questions, suggestions or improvements? Feel free to open issues or PRs!
