from __future__ import annotations

import argparse
import json
import platform
import shutil
import subprocess
import sys
from pathlib import Path


PROJECT_BASE = Path(__file__).parent.parent.resolve()


def detect_os() -> str:
    system = platform.system().lower()
    return {
        "windows": "win",
        "darwin": "macos",
        "linux": "linux",
    }.get(system, system)


def detect_arch() -> str:
    machine = platform.machine().lower()
    return {
        "amd64": "x86_64",
        "x86_64": "x86_64",
        "arm64": "aarch64",
        "aarch64": "aarch64",
    }.get(machine, machine)


def read_version_from_interface(install_dir: Path) -> str:
    candidates = [
        install_dir / "interface.json",
        PROJECT_BASE / "assets" / "interface.json",
    ]

    for candidate in candidates:
        if not candidate.is_file():
            continue
        with open(candidate, "r", encoding="utf-8") as file:
            data = json.load(file)
        version = data.get("version")
        if isinstance(version, str) and version.strip():
            return version.strip()

    raise ValueError("Unable to determine package version from interface.json")


def validate_install_dir(install_dir: Path, os_name: str) -> None:
    launcher_name = "MDA.exe" if os_name == "win" else "MDA"
    agent_name = "go-service.exe" if os_name == "win" else "go-service"

    required_paths = [
        install_dir / launcher_name,
        install_dir / "agent" / agent_name,
        install_dir / "maafw",
        install_dir / "resource",
        install_dir / "tasks",
        install_dir / "locales",
        install_dir / "interface.json",
    ]

    missing_paths = [path for path in required_paths if not path.exists()]
    if missing_paths:
        missing_text = "\n".join(str(path) for path in missing_paths)
        raise FileNotFoundError(f"Install directory is incomplete:\n{missing_text}")


def run_build(py_executable: str, version: str | None) -> None:
    command = [py_executable, str(PROJECT_BASE / "tools" / "setup_workspace.py"), "--ci"]
    if version:
        print(
            "Packaging uses the current install directory version; "
            f"explicit package version override is {version}."
        )

    print(f"Running build command: {' '.join(command)}")
    subprocess.run(command, cwd=PROJECT_BASE, check=True)


def create_archive(install_dir: Path, output_dir: Path, archive_name: str) -> Path:
    output_dir.mkdir(parents=True, exist_ok=True)
    archive_base = output_dir / archive_name
    archive_path = archive_base.with_suffix(".zip")

    if archive_path.exists():
        archive_path.unlink()

    created_path = shutil.make_archive(
        str(archive_base),
        "zip",
        root_dir=install_dir,
        base_dir=".",
    )
    return Path(created_path)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Package install directory into a release zip")
    parser.add_argument("--build", action="store_true", help="Run tools/setup_workspace.py --ci before packaging")
    parser.add_argument("--version", help="Package version override, defaults to install/interface.json")
    parser.add_argument("--os", default=detect_os(), help="Package OS name, defaults to current platform")
    parser.add_argument("--arch", default=detect_arch(), help="Package architecture, defaults to current platform")
    parser.add_argument("--install-dir", default=str(PROJECT_BASE / "install"), help="Install directory to package")
    parser.add_argument("--output-dir", default=str(PROJECT_BASE / "dist"), help="Output directory for the zip")
    return parser.parse_args()


def main() -> None:
    args = parse_args()
    install_dir = Path(args.install_dir).resolve()
    output_dir = Path(args.output_dir).resolve()

    py_executable = sys.executable

    if args.build:
        run_build(py_executable, args.version)

    validate_install_dir(install_dir, args.os)

    version = args.version or read_version_from_interface(install_dir)
    archive_name = f"MDA-{args.os}-{args.arch}-{version}"
    archive_path = create_archive(install_dir, output_dir, archive_name)

    print(f"Release package created: {archive_path}")


if __name__ == "__main__":
    main()
