#!/usr/bin/env python3
"""Workspace-local Antigravity MCP configuration generator.

Reads and writes <workspace>/.agents/mcp_config.json, registering one
local stdio server while preserving unrelated servers and settings.
Global user configuration is never touched.

Exit codes: 0 success (including idempotent no-op); 2 usage or
configuration errors; 1 unexpected failures.
"""

import argparse
import json
import os
import shutil
import sys
import time

SERVER_NAME_DEFAULT = "phystrace"


class ConfigError(Exception):
    """Actionable configuration failure (exit 2)."""


def config_path(workspace):
    return os.path.join(workspace, ".agents", "mcp_config.json")


def load_config(path):
    """Load config, tolerating a missing or empty file (fresh document).

    Invalid JSON, a non-object root, or a non-object mcpServers mapping
    are errors: silently repairing them could discard user settings.
    """
    try:
        with open(path, "r", encoding="utf-8") as f:
            text = f.read()
    except FileNotFoundError:
        return {"mcpServers": {}}
    if not text.strip():
        return {"mcpServers": {}}
    try:
        data = json.loads(text)
    except json.JSONDecodeError as e:
        raise ConfigError("existing config is not valid JSON: %s" % e)
    if not isinstance(data, dict):
        raise ConfigError("existing config root is not an object")
    servers = data.get("mcpServers", {})
    if not isinstance(servers, dict):
        raise ConfigError("mcpServers is not an object")
    data["mcpServers"] = servers
    return data


def validate_command(command):
    if not os.path.isabs(command):
        raise ConfigError("server command must be an absolute path: %r" % command)
    if not os.path.isfile(command):
        raise ConfigError("server command does not exist: %s" % command)
    if not os.access(command, os.X_OK):
        raise ConfigError("server command is not executable: %s" % command)


def ensure_server(config, name, command, workspace):
    """Set servers[name] to the stdio entry with an explicit workspace
    binding. Returns True when changed.

    The entry carries args ["--workspace", <absolute workspace>] so the
    server does not depend on the client's working directory. Everything
    else in the document is preserved semantically. Note the document is
    reserialized on write, so exact original bytes/layout are not
    preserved -- values are.
    """
    entry = {"command": command, "args": ["--workspace", workspace]}
    if config["mcpServers"].get(name) == entry:
        return False
    config["mcpServers"][name] = entry
    return True


def render(config):
    return json.dumps(config, indent=4, sort_keys=True) + "\n"


def write_config(workspace, config):
    """Write config atomically; back up the previous file only on change.

    Returns (changed, backup_path_or_None). A repeat run with identical
    content is a no-op: no backup, no rewrite.
    """
    path = config_path(workspace)
    new_text = render(config)
    try:
        with open(path, "r", encoding="utf-8") as f:
            old_text = f.read()
    except FileNotFoundError:
        old_text = None
    if old_text == new_text:
        return False, None
    # Re-validate what we are about to write.
    json.loads(new_text)
    backup = None
    if old_text is not None:
        stamp = time.strftime("%Y%m%dT%H%M%S")
        backup = "%s.bak.%s" % (path, stamp)
        shutil.copy2(path, backup)
    agents_dir = os.path.dirname(path)
    os.makedirs(agents_dir, exist_ok=True)
    tmp_path = "%s.tmp.%d" % (path, os.getpid())
    with open(tmp_path, "w", encoding="utf-8") as f:
        f.write(new_text)
    # Atomic replacement on POSIX. On platforms where an existing
    # destination cannot be replaced atomically, os.replace raises and
    # the original is left untouched (verified by test, not retried:
    # retrying a failed replace risks partial state).
    os.replace(tmp_path, path)
    return True, backup


def main(argv):
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--workspace", required=True, help="demo-workspace root")
    parser.add_argument("--binary", required=True, help="absolute path to the MCP server binary")
    parser.add_argument("--server-name", default=SERVER_NAME_DEFAULT)
    args = parser.parse_args(argv)
    try:
        if not os.path.isdir(args.workspace):
            raise ConfigError("workspace is not a directory: %s" % args.workspace)
        workspace = os.path.realpath(os.path.abspath(args.workspace))
        validate_command(args.binary)
        path = config_path(workspace)
        config = load_config(path)
        changed = ensure_server(config, args.server_name, args.binary, workspace)
        backup = None
        if changed:
            changed, backup = write_config(workspace, config)
        print(json.dumps({"config": path, "changed": changed, "backup": backup}))
        return 0
    except ConfigError as e:
        print("antigravity_config: error: %s" % e, file=sys.stderr)
        return 2
    except OSError as e:
        print("antigravity_config: filesystem error: %s" % e, file=sys.stderr)
        return 1


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
