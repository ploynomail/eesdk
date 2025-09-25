import ctypes
import json
from typing import Optional, Dict, Any
from ctypes import Structure, c_char_p, c_int, c_size_t, POINTER

class EESdkConfig(Structure):
    _fields_ = [
        ("base_url", c_char_p),
        ("api_key", c_char_p),
        ("timeout_seconds", c_int),
        ("retry_count", c_int),
    ]

class EESdkResponse(Structure):
    _fields_ = [
        ("status_code", c_int),
        ("body", c_char_p),
        ("body_length", c_size_t),
        ("error_message", c_char_p),
    ]

class EESdkClient:
    def __init__(self, base_url: str, api_key: str, timeout: int = 30, retry_count: int = 3):
        # Load the shared library
        self._lib = ctypes.CDLL("../../build/lib/libee_sdk.so")
        
        # Define function signatures
        self._lib.ee_client_create.argtypes = [POINTER(EESdkConfig)]
        self._lib.ee_client_create.restype = ctypes.c_void_p
        
        self._lib.ee_client_destroy.argtypes = [ctypes.c_void_p]
        self._lib.ee_client_destroy.restype = None
        
        self._lib.ee_execute_request.argtypes = [ctypes.c_void_p, c_char_p, c_char_p, c_char_p]
        self._lib.ee_execute_request.restype = POINTER(EESdkResponse)
        
        self._lib.ee_health_check.argtypes = [ctypes.c_void_p]
        self._lib.ee_health_check.restype = c_int
        
        self._lib.ee_response_free.argtypes = [POINTER(EESdkResponse)]
        self._lib.ee_response_free.restype = None
        
        # Create client
        config = EESdkConfig(
            base_url=base_url.encode('utf-8'),
            api_key=api_key.encode('utf-8'),
            timeout_seconds=timeout,
            retry_count=retry_count
        )
        
        self._client = self._lib.ee_client_create(ctypes.byref(config))
        if not self._client:
            raise RuntimeError("Failed to create SDK client")
    
    def __del__(self):
        if hasattr(self, '_client') and self._client:
            self._lib.ee_client_destroy(self._client)
    
    def execute_request(self, method: str, path: str, body: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        body_str = json.dumps(body).encode('utf-8') if body else b""
        
        response_ptr = self._lib.ee_execute_request(
            self._client,
            method.encode('utf-8'),
            path.encode('utf-8'),
            body_str
        )
        
        if not response_ptr:
            raise RuntimeError("Request execution failed")
        
        try:
            response = response_ptr.contents
            
            if response.error_message:
                error_msg = response.error_message.decode('utf-8')
                raise RuntimeError(f"Request failed: {error_msg}")
            
            result = {
                'status_code': response.status_code,
                'body': response.body.decode('utf-8') if response.body else "",
            }
            
            # Try to parse JSON response
            if result['body']:
                try:
                    result['body'] = json.loads(result['body'])
                except json.JSONDecodeError:
                    pass  # Keep as string if not valid JSON
            
            return result
        
        finally:
            self._lib.ee_response_free(response_ptr)
    
    def health_check(self) -> bool:
        result = self._lib.ee_health_check(self._client)
        return result == 0
