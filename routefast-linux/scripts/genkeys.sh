#!/usr/bin/env bash
set -eu
mkdir -p configs/keys
python3 - <<'PY'
import base64
try:
    from cryptography.hazmat.primitives.asymmetric.ed25519 import Ed25519PrivateKey
    from cryptography.hazmat.primitives.serialization import Encoding, PrivateFormat, NoEncryption, PublicFormat
except Exception as e:
    raise SystemExit("Python package 'cryptography' is required. Install it with: pip install cryptography") from e
for node in ["rt-a", "rt-b", "rt-c"]:
    priv = Ed25519PrivateKey.generate()
    pub = priv.public_key()
    priv_raw = priv.private_bytes(Encoding.Raw, PrivateFormat.Raw, NoEncryption())
    pub_raw = pub.public_bytes(Encoding.Raw, PublicFormat.Raw)
    open(f"configs/keys/{node}.key", "w").write(base64.b64encode(priv_raw).decode() + "\n")
    print(node, "pubkey =", "ed25519:" + base64.b64encode(pub_raw).decode())
PY
