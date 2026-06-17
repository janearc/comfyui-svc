import urllib.request
import urllib.error
import json
import sys

def test_traefik_routing():
    print("Testing Traefik routing for comfyui.localhost...")
    req = urllib.request.Request("http://127.0.0.1/", headers={"Host": "comfyui.localhost"})
    try:
        with urllib.request.urlopen(req) as response:
            assert response.status == 200, f"Expected status 200, got {response.status}"
            print("✓ Root routing successful.")
    except urllib.error.URLError as e:
        print(f"✗ Failed to connect or route: {e}")
        sys.exit(1)

def test_health_endpoint():
    print("Testing /health endpoint on sidecar...")
    req = urllib.request.Request("http://127.0.0.1/health", headers={"Host": "comfyui.localhost"})
    try:
        with urllib.request.urlopen(req) as response:
            assert response.status == 200, f"Expected status 200, got {response.status}"
            data = json.loads(response.read().decode('utf-8'))
            assert data.get("status") == "ok", f"Expected status 'ok', got {data.get('status')}"
            print("✓ Health endpoint returned healthy status.")
    except urllib.error.URLError as e:
        print(f"✗ Failed to connect or fetch health: {e}")
        sys.exit(1)

def test_metal_backend():
    print("Testing /system_stats to ensure backend is using Metal (MPS)...")
    req = urllib.request.Request("http://127.0.0.1/system_stats", headers={"Host": "comfyui.localhost"})
    try:
        with urllib.request.urlopen(req) as response:
            assert response.status == 200, f"Expected status 200, got {response.status}"
            data = json.loads(response.read().decode('utf-8'))
            devices = data.get("devices", [])
            assert any(d.get("type") == "mps" for d in devices), f"Expected 'mps' device, but got: {devices}"
            print("✓ Backend is confirmed to be using Metal (MPS).")
    except urllib.error.URLError as e:
        print(f"✗ Failed to connect or fetch system stats: {e}")
        sys.exit(1)

if __name__ == "__main__":
    try:
        test_traefik_routing()
        test_health_endpoint()
        test_metal_backend()
        print("All regression tests passed successfully!")
    except AssertionError as ae:
        print(f"Assertion Error: {ae}")
        sys.exit(1)
