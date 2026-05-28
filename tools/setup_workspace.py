import argparse
import sys
import shutil
import subprocess
import platform
from pathlib import Path

from cli_support import Console, init_localization


PROJECT_BASE: Path = Path(__file__).parent.parent.resolve()


def create_directory_link(src: Path, dst: Path) -> bool:
    """
    在指定位置创建一个指定目录的链接
    - Windows：Junction
    - Unix/macOS：symlink
    """
    if dst.exists() or dst.is_symlink():
        if dst.is_dir() and not dst.is_symlink():
            try:
                dst.rmdir()
            except OSError:
                shutil.rmtree(dst)
        else:
            dst.unlink(missing_ok=True)

    dst.parent.mkdir(parents=True, exist_ok=True)

    if platform.system() == "Windows":
        result = subprocess.run(
            ["cmd", "/c", "mklink", "/J", str(dst), str(src)],
            capture_output=True,
            text=True,
        )
        if result.returncode != 0:
            print(Console.err(t("err_create_junction_failed", stderr=result.stderr)))
            return False
    else:
        dst.symlink_to(src)

    return True


LOCALS_DIR = Path(__file__).parent / "locals" / "setup_workspace"


_local_t = lambda key, **kwargs: key.format(**kwargs) if kwargs else key


def init_local() -> None:
    global _local_t
    t_func, load_error_path = init_localization(LOCALS_DIR)
    _local_t = t_func
    if load_error_path:
        print(Console.err(t("error_load_locale", path=load_error_path)))


def t(key: str, **kwargs) -> str:
    return _local_t(key, **kwargs)


try:
    OS_KEYWORD: str = {
        "windows": "win",
        "linux": "linux",
        "darwin": "macos",
    }[platform.system().lower()]
except KeyError as e:
    raise RuntimeError(
        f"Unrecognized operating system: {platform.system().lower()}"
    ) from e

try:
    ARCH_KEYWORD: str = {
        "amd64": "x86_64",
        "x86_64": "x86_64",
        "aarch64": "aarch64",
        "arm64": "aarch64",
    }[platform.machine().lower()]
except KeyError as e:
    raise RuntimeError(
        f"Unrecognized architecture: {platform.machine().lower()}"
    ) from e

try:
    MFW_DIST_NAME: str = {
        "win": "MaaFramework.dll",
        "linux": "libMaaFramework.so",
        "macos": "libMaaFramework.dylib",
    }[OS_KEYWORD]
except KeyError as e:
    raise RuntimeError(f"Unsupported OS for MaaFramework: {OS_KEYWORD}") from e

MXU_DIST_NAME: str = "mxu.exe" if OS_KEYWORD == "win" else "mxu"
MXU_LAUNCHER_NAME: str = "MDA.exe" if OS_KEYWORD == "win" else "MDA"
CACHE_DIR: Path = PROJECT_BASE / ".cache"


def configure_token() -> None:
    # 自动更新逻辑已移除，保留函数名仅为兼容旧调用面。
    return None


def run_command(
    cmd: list[str] | str, cwd: Path | str | None = None, shell: bool = False
) -> bool:
    cmd_str = " ".join(cmd) if isinstance(cmd, list) else str(cmd)
    print(f"{Console.info(t('cmd_prefix'))} {cmd_str}")
    try:
        subprocess.check_call(cmd, cwd=cwd or PROJECT_BASE, shell=shell)
        print(Console.ok(t("inf_command_success", cmd=cmd_str)))
        return True
    except subprocess.CalledProcessError as e:
        print(Console.err(t("err_command_failed", cmd=cmd_str, error=e)))
        return False


def update_submodules(skip_if_exist: bool = True) -> bool:
    print(Console.hdr(t("inf_check_submodules")))

    common_assets_path = PROJECT_BASE / "assets" / "MaaCommonAssets"
    if skip_if_exist and common_assets_path.exists() and any(common_assets_path.iterdir()):
        print(Console.ok(t("inf_submodules_exist")))
        return True

    print(Console.info(t("inf_updating_submodules")))
    return run_command(["git", "submodule", "update", "--init", "--recursive"])


def run_build_script(ci_mode: bool = False) -> bool:
    """执行 build_and_install.py"""
    print(Console.hdr(t("inf_run_build_script")))
    cmd = [sys.executable, str(PROJECT_BASE / "tools" / "build_and_install.py")]
    if ci_mode:
        cmd.append("--ci")
    return run_command(cmd)


def clean_cache() -> None:
    if not CACHE_DIR.exists():
        print(Console.info(t("inf_cache_empty")))
        return
    total_size = 0
    count = 0
    for f in CACHE_DIR.iterdir():
        if f.is_file():
            total_size += f.stat().st_size
            count += 1
    if count == 0:
        print(Console.info(t("inf_cache_empty")))
        return
    size_mb = total_size / (1024 * 1024)
    print(Console.info(t("inf_cache_summary", count=count, size=f"{size_mb:.1f} MB")))
    try:
        shutil.rmtree(CACHE_DIR)
        print(Console.ok(t("inf_cache_purged")))
    except OSError as e:
        print(Console.warn(t("wrn_cache_clean_failed", path=CACHE_DIR, error=e)))


def install_maafw(
    install_root: Path,
    skip_if_exist: bool = True,
    update_mode: bool = False,
    local_version: str | None = None,
) -> tuple[bool, str | None, bool]:
    _ = skip_if_exist, update_mode
    real_install_root = install_root.resolve()
    maafw_dest = real_install_root / "maafw"
    maafw_deps = PROJECT_BASE / "deps"
    maafw_installed = maafw_deps.exists() and any(maafw_deps.iterdir())

    # 自动更新已移除：仅使用本地依赖，不再访问 GitHub Release。
    # skip_if_exist / update_mode 参数保留以兼容原调用签名。
    if maafw_installed:
        print(Console.ok(t("inf_maafw_installed_skip")))
        return True, local_version or "local", False
    if maafw_dest.exists():
        return True, local_version or "local", False
    print(Console.err(t("err_maafw_url_not_found")))
    print(Console.err("[ERR] Auto update disabled: MaaFramework local dependency not found."))
    return False, local_version, False


def install_mxu(
    install_root: Path,
    skip_if_exist: bool = True,
    update_mode: bool = False,
    local_version: str | None = None,
) -> tuple[bool, str | None, bool]:
    _ = skip_if_exist, update_mode
    real_install_root = install_root.resolve()
    mxu_path = real_install_root / MXU_DIST_NAME
    launcher_path = real_install_root / MXU_LAUNCHER_NAME
    launcher_installed = launcher_path.exists()

    # 自动更新已移除：仅使用本地安装目录，不再访问 GitHub Release。
    # skip_if_exist / update_mode 参数保留以兼容原调用签名。
    if launcher_installed:
        print(Console.ok(t("inf_mxu_installed_skip")))
        return True, local_version or "local", False
    if mxu_path.exists():
        return True, local_version or "local", False
    print(Console.err(t("err_mxu_url_not_found")))
    print(Console.err("[ERR] Auto update disabled: MXU local executable not found."))
    return False, local_version, False


def main() -> None:
    init_local()

    parser = argparse.ArgumentParser(description=t("description"))
    parser.add_argument("--update", action="store_true", help=t("arg_update"))
    parser.add_argument("--ci", action="store_true", help=t("arg_ci"))
    parser.add_argument("--clean-cache", action="store_true", help=t("arg_clean_cache"))
    args = parser.parse_args()

    if args.clean_cache:
        clean_cache()
        return

    install_dir = PROJECT_BASE / "install"

    print(Console.hdr(t("header_workspace_init")))

    # 1. Update submodules
    if not update_submodules(skip_if_exist=not args.update):
        print(Console.err(t("fatal_submodule_failed")))
        sys.exit(1)

    # 2. Build and install (delegated to build_and_install.py)
    print(Console.hdr(t("header_build_and_install")))
    if not run_build_script(ci_mode=args.ci):
        print(Console.err(t("fatal_build_failed")))
        sys.exit(1)

    # 3. Download MaaFramework & MXU
    print(Console.hdr(t("header_download_deps")))
    maafw_ok, _, _ = install_maafw(
        install_dir,
        skip_if_exist=not args.update,
        update_mode=args.update,
        local_version=None,
    )
    if not maafw_ok:
        print(Console.err(t("fatal_maafw_failed")))
        sys.exit(1)

    mxu_ok, _, _ = install_mxu(
        install_dir,
        skip_if_exist=not args.update,
        update_mode=args.update,
        local_version=None,
    )
    if not mxu_ok:
        print(Console.err(t("fatal_mxu_failed")))
        sys.exit(1)

    print(Console.ok(t("header_setup_complete")))
    print(Console.info(t("inf_workspace_ready", mxu_path=install_dir / MXU_LAUNCHER_NAME)))
    print(Console.info(t("inf_install_dir_hint", install_dir=install_dir)))


if __name__ == "__main__":
    main()
