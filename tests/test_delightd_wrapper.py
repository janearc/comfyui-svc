import os
import sys

def test_delightd_bash_wrapper():
    print("Testing if delightd bash wrapper exists for comfyui...")
    wrapper_path = os.path.expanduser("~/var/bin/comfyui")
    
    assert os.path.exists(wrapper_path), f"Wrapper not found at {wrapper_path}"
    assert os.access(wrapper_path, os.X_OK), f"Wrapper at {wrapper_path} is not executable"
    
    with open(wrapper_path, "r") as f:
        content = f.read()
        assert "docker exec" in content, "Wrapper does not contain 'docker exec' command"
        assert "comfyui-comfyui-1" in content, "Wrapper does not target the comfyui container"
        
    print("✓ Delightd bash wrapper exists and is correctly configured.")

if __name__ == "__main__":
    try:
        test_delightd_bash_wrapper()
    except AssertionError as ae:
        print(f"Assertion Error: {ae}")
        sys.exit(1)
