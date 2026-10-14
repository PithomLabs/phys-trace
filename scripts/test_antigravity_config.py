#!/usr/bin/env python3
"""Unit tests for antigravity_config.py (stdlib unittest, no pytest needed).

Run:  python3 -m unittest -v test_antigravity_config
from the scripts/ directory.
"""

import json
import os
import stat
import subprocess
import sys
import tempfile
import unittest

import antigravity_config as ac


class ConfigTest(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        self.ws = os.path.join(self.tmp.name, "ws")
        os.makedirs(self.ws)
        self.bin = os.path.join(self.tmp.name, "phystrace-mcp")
        with open(self.bin, "w") as f:
            f.write("#!/bin/sh\nexit 0\n")
        os.chmod(self.bin, 0o755)

    def test_fresh_write(self):
        config = ac.load_config(ac.config_path(self.ws))
        self.assertTrue(ac.ensure_server(config, "phystrace", self.bin))
        changed, backup = ac.write_config(self.ws, config)
        self.assertTrue(changed)
        self.assertIsNone(backup)
        with open(ac.config_path(self.ws)) as f:
            data = json.load(f)
        self.assertEqual(data["mcpServers"]["phystrace"]["command"], self.bin)

    def test_idempotent_no_backup(self):
        config = ac.load_config(ac.config_path(self.ws))
        ac.ensure_server(config, "phystrace", self.bin)
        ac.write_config(self.ws, config)
        config2 = ac.load_config(ac.config_path(self.ws))
        self.assertFalse(ac.ensure_server(config2, "phystrace", self.bin))
        changed, backup = ac.write_config(self.ws, config2)
        self.assertFalse(changed)
        self.assertIsNone(backup)

    def test_preserves_unrelated_servers(self):
        path = ac.config_path(self.ws)
        os.makedirs(os.path.dirname(path), exist_ok=True)
        with open(path, "w") as f:
            json.dump({"mcpServers": {"other": {"command": "/usr/bin/other"}},
                       "custom": [1, 2]}, f)
        config = ac.load_config(path)
        self.assertTrue(ac.ensure_server(config, "phystrace", self.bin))
        changed, backup = ac.write_config(self.ws, config)
        self.assertTrue(changed)
        self.assertIsNotNone(backup)
        with open(path) as f:
            data = json.load(f)
        # Values preserved (reserialization may reformat bytes).
        self.assertEqual(data["mcpServers"]["other"], {"command": "/usr/bin/other"})
        self.assertEqual(data["custom"], [1, 2])
        self.assertEqual(data["mcpServers"]["phystrace"]["command"], self.bin)

    def test_empty_file_is_fresh(self):
        path = ac.config_path(self.ws)
        os.makedirs(os.path.dirname(path), exist_ok=True)
        with open(path, "w") as f:
            f.write("")
        config = ac.load_config(path)
        self.assertEqual(config, {"mcpServers": {}})

    def test_invalid_json_rejected(self):
        path = ac.config_path(self.ws)
        os.makedirs(os.path.dirname(path), exist_ok=True)
        with open(path, "w") as f:
            f.write("{not json")
        with self.assertRaises(ac.ConfigError):
            ac.load_config(path)

    def test_command_validation(self):
        with self.assertRaises(ac.ConfigError):
            ac.validate_command("relative/bin")
        with self.assertRaises(ac.ConfigError):
            ac.validate_command("/nonexistent-binary-xyz")
        plain = os.path.join(self.tmp.name, "plain")
        with open(plain, "w") as f:
            f.write("x")
        with self.assertRaises(ac.ConfigError):
            ac.validate_command(plain)
        ac.validate_command(self.bin)

    def test_failed_write_preserves_original(self):
        path = ac.config_path(self.ws)
        os.makedirs(os.path.dirname(path), exist_ok=True)
        with open(path, "w") as f:
            json.dump({"mcpServers": {"other": {"command": "/usr/bin/other"}}}, f)
        # Make the destination directory read-only so replacement fails.
        os.chmod(os.path.dirname(path), 0o555)
        self.addCleanup(os.chmod, os.path.dirname(path), 0o755)
        config = ac.load_config(path)
        ac.ensure_server(config, "phystrace", self.bin)
        with self.assertRaises(OSError):
            ac.write_config(self.ws, config)
        with open(path) as f:
            self.assertEqual(json.load(f)["mcpServers"], {"other": {"command": "/usr/bin/other"}})

    def test_cli_end_to_end(self):
        r = subprocess.run(
            [sys.executable, "antigravity_config.py", "--workspace", self.ws,
             "--binary", self.bin],
            capture_output=True, text=True, cwd=os.path.dirname(os.path.abspath(__file__)))
        self.assertEqual(r.returncode, 0, r.stderr)
        summary = json.loads(r.stdout)
        self.assertTrue(summary["changed"])
        r2 = subprocess.run(
            [sys.executable, "antigravity_config.py", "--workspace", self.ws,
             "--binary", self.bin],
            capture_output=True, text=True, cwd=os.path.dirname(os.path.abspath(__file__)))
        self.assertEqual(r2.returncode, 0, r2.stderr)
        self.assertFalse(json.loads(r2.stdout)["changed"])

    def test_cli_bad_binary_exit_2(self):
        r = subprocess.run(
            [sys.executable, "antigravity_config.py", "--workspace", self.ws,
             "--binary", "/nonexistent-xyz"],
            capture_output=True, text=True, cwd=os.path.dirname(os.path.abspath(__file__)))
        self.assertEqual(r.returncode, 2)


if __name__ == "__main__":
    unittest.main()
