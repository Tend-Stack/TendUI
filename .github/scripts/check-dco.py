#!/usr/bin/env python3
"""Check author-matching DCO trailers without executing contribution code."""
import re
import subprocess
import sys


def signed_off(name, email, message):
    trailers = subprocess.run(
        ["git", "interpret-trailers", "--parse"], input=message,
        text=True, capture_output=True, check=True,
    ).stdout
    for line in trailers.splitlines():
        match = re.fullmatch(r"Signed-off-by:\s*(.+?)\s*<([^<>\s]+@[^<>\s]+)>", line, re.I)
        if match and match[1].strip() == name.strip() and match[2].casefold() == email.casefold():
            return True
    return False


def main():
    if len(sys.argv) != 3 or not all(re.fullmatch(r"[0-9a-f]{40}", s) for s in sys.argv[1:]):
        raise SystemExit("Usage: check-dco.py BASE_SHA HEAD_SHA (full commit hashes)")
    commits = subprocess.check_output(
        ["git", "rev-list", "--no-merges", sys.argv[1] + ".." + sys.argv[2]], text=True,
    ).splitlines()
    failed = []
    for commit in commits:
        author, email, message = subprocess.check_output(
            ["git", "show", "-s", "--format=%an%x00%ae%x00%B", commit], text=True,
        ).split("\0", 2)
        if not signed_off(author, email, message):
            failed.append(commit[:12])
    if failed:
        print("Missing author-matching DCO sign-off in: " + ", ".join(failed))
        print("Read CONTRIBUTING.md; sign only your own contributions with git commit -s.")
        return 1
    print(f"DCO: {len(commits)} non-merge commits have author-matching sign-offs.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
