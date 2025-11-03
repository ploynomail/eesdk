import ctypes
from typing import Optional, List
from ctypes import Structure, c_char_p, c_int, c_uint, c_size_t, POINTER, c_void_p

class EESdkConfig(Structure):
    _fields_ = [
        ("server_addr", c_char_p),
        ("access_id", c_char_p),
        ("access_key", c_char_p),
    ]

class EESubject(Structure):
    _fields_ = [
        ("common_name", c_char_p),
        ("country", POINTER(c_char_p)),
        ("country_len", c_int),
        ("organization", POINTER(c_char_p)),
        ("organization_len", c_int),
        ("organizational_unit", POINTER(c_char_p)),
        ("organizational_unit_len", c_int),
        ("locality", POINTER(c_char_p)),
        ("locality_len", c_int),
        ("province", POINTER(c_char_p)),
        ("province_len", c_int),
    ]

class EECertRequest(Structure):
    _fields_ = [
        ("public_key", POINTER(ctypes.c_ubyte)),
        ("public_key_len", c_size_t),
        ("kem_public_key", POINTER(ctypes.c_ubyte)),
        ("kem_public_key_len", c_size_t),
        ("ca_cert_id", c_uint),
        ("template_id", c_uint),
        ("subject", EESubject),
    ]

class EECertResponse(Structure):
    _fields_ = [
        ("signer_cert", POINTER(ctypes.c_ubyte)),
        ("signer_cert_len", c_size_t),
        ("enc_cert", POINTER(ctypes.c_ubyte)),
        ("enc_cert_len", c_size_t),
        ("encrypted_private_key", POINTER(ctypes.c_ubyte)),
        ("encrypted_private_key_len", c_size_t),
        ("error_message", c_char_p),
    ]

class EERevokeRequest(Structure):
    _fields_ = [
        ("serial_number", c_char_p),
        ("reason", c_char_p),
    ]

class EEResponse(Structure):
    _fields_ = [
        ("success", c_int),
        ("error_message", c_char_p),
    ]

class Subject:
    def __init__(self, common_name: str, country: Optional[List[str]] = None,
                 organization: Optional[List[str]] = None,
                 organizational_unit: Optional[List[str]] = None,
                 locality: Optional[List[str]] = None,
                 province: Optional[List[str]] = None):
        self.common_name = common_name
        self.country = country or []
        self.organization = organization or []
        self.organizational_unit = organizational_unit or []
        self.locality = locality or []
        self.province = province or []

class CertificateResponse:
    def __init__(self, signer_cert: bytes, enc_cert: bytes, encrypted_private_key: bytes):
        self.signer_cert = signer_cert
        self.enc_cert = enc_cert
        self.encrypted_private_key = encrypted_private_key

class EESdkClient:
    def __init__(self, server_addr: str, access_id: str, access_key: str):
        # 加载共享库
        self._lib = ctypes.CDLL("libee_sdk.so")
        
        # 定义函数签名
        self._lib.ee_client_create.argtypes = [POINTER(EESdkConfig)]
        self._lib.ee_client_create.restype = c_void_p
        
        self._lib.ee_client_destroy.argtypes = [c_void_p]
        self._lib.ee_client_destroy.restype = None
        
        self._lib.ee_request_certificate.argtypes = [c_void_p, POINTER(EECertRequest)]
        self._lib.ee_request_certificate.restype = POINTER(EECertResponse)
        
        self._lib.ee_request_certificate_pkcs10.argtypes = [c_void_p, POINTER(ctypes.c_ubyte), c_size_t]
        self._lib.ee_request_certificate_pkcs10.restype = POINTER(EECertResponse)
        
        self._lib.ee_revoke_certificate.argtypes = [c_void_p, POINTER(EERevokeRequest)]
        self._lib.ee_revoke_certificate.restype = POINTER(EEResponse)
        
        self._lib.ee_confirm_certificate.argtypes = [c_void_p, POINTER(ctypes.c_ubyte), c_size_t, c_int]
        self._lib.ee_confirm_certificate.restype = POINTER(EEResponse)
        
        self._lib.ee_cert_response_free.argtypes = [POINTER(EECertResponse)]
        self._lib.ee_cert_response_free.restype = None
        
        self._lib.ee_response_free.argtypes = [POINTER(EEResponse)]
        self._lib.ee_response_free.restype = None
        
        # 创建客户端
        config = EESdkConfig(
            server_addr=server_addr.encode('utf-8'),
            access_id=access_id.encode('utf-8'),
            access_key=access_key.encode('utf-8')
        )
        
        self._client = self._lib.ee_client_create(ctypes.byref(config))
        if not self._client:
            raise RuntimeError("Failed to create SDK client")
    
    def __del__(self):
        if hasattr(self, '_client') and self._client:
            self._lib.ee_client_destroy(self._client)
    
    def request_certificate(self, public_key: bytes, kem_public_key: bytes,
                          ca_cert_id: int, template_id: int,
                          subject: Subject) -> CertificateResponse:
        # 准备主题
        c_subject = EESubject()
        c_subject.common_name = subject.common_name.encode('utf-8')
        
        # 转换字符串列表
        country_arr = (c_char_p * len(subject.country))(*[s.encode('utf-8') for s in subject.country])
        c_subject.country = country_arr
        c_subject.country_len = len(subject.country)
        
        org_arr = (c_char_p * len(subject.organization))(*[s.encode('utf-8') for s in subject.organization])
        c_subject.organization = org_arr
        c_subject.organization_len = len(subject.organization)
        
        # 准备请求
        req = EECertRequest()
        req.public_key = (ctypes.c_ubyte * len(public_key))(*public_key)
        req.public_key_len = len(public_key)
        req.kem_public_key = (ctypes.c_ubyte * len(kem_public_key))(*kem_public_key)
        req.kem_public_key_len = len(kem_public_key)
        req.ca_cert_id = ca_cert_id
        req.template_id = template_id
        req.subject = c_subject
        
        response_ptr = self._lib.ee_request_certificate(self._client, ctypes.byref(req))
        
        if not response_ptr:
            raise RuntimeError("Request execution failed")
        
        try:
            response = response_ptr.contents
            
            if response.error_message:
                error_msg = response.error_message.decode('utf-8')
                raise RuntimeError(f"Request failed: {error_msg}")
            
            # 提取响应数据
            signer_cert = bytes(ctypes.cast(response.signer_cert, POINTER(ctypes.c_ubyte * response.signer_cert_len)).contents)
            enc_cert = bytes(ctypes.cast(response.enc_cert, POINTER(ctypes.c_ubyte * response.enc_cert_len)).contents)
            encrypted_key = bytes(ctypes.cast(response.encrypted_private_key, POINTER(ctypes.c_ubyte * response.encrypted_private_key_len)).contents)
            
            return CertificateResponse(signer_cert, enc_cert, encrypted_key)
        
        finally:
            self._lib.ee_cert_response_free(response_ptr)
    
    def request_certificate_with_pkcs10(self, csr: bytes) -> CertificateResponse:
        csr_array = (ctypes.c_ubyte * len(csr))(*csr)
        response_ptr = self._lib.ee_request_certificate_pkcs10(
            self._client, csr_array, len(csr)
        )
        
        if not response_ptr:
            raise RuntimeError("Request execution failed")
        
        try:
            response = response_ptr.contents
            
            if response.error_message:
                error_msg = response.error_message.decode('utf-8')
                raise RuntimeError(f"Request failed: {error_msg}")
            
            signer_cert = bytes(ctypes.cast(response.signer_cert, POINTER(ctypes.c_ubyte * response.signer_cert_len)).contents)
            enc_cert = bytes(ctypes.cast(response.enc_cert, POINTER(ctypes.c_ubyte * response.enc_cert_len)).contents)
            
            return CertificateResponse(signer_cert, enc_cert, b'')
        
        finally:
            self._lib.ee_cert_response_free(response_ptr)
    
    def revoke_certificate(self, serial_number: str, reason: str) -> bool:
        req = EERevokeRequest(
            serial_number=serial_number.encode('utf-8'),
            reason=reason.encode('utf-8')
        )
        
        response_ptr = self._lib.ee_revoke_certificate(self._client, ctypes.byref(req))
        
        if not response_ptr:
            raise RuntimeError("Revoke execution failed")
        
        try:
            response = response_ptr.contents
            
            if response.error_message:
                error_msg = response.error_message.decode('utf-8')
                raise RuntimeError(f"Revoke failed: {error_msg}")
            
            return response.success == 1
        
        finally:
            self._lib.ee_response_free(response_ptr)
    
    def confirm_certificate(self, cert_hash: bytes, cert_req_id: int) -> bool:
        hash_array = (ctypes.c_ubyte * len(cert_hash))(*cert_hash)
        response_ptr = self._lib.ee_confirm_certificate(
            self._client, hash_array, len(cert_hash), cert_req_id
        )
        
        if not response_ptr:
            raise RuntimeError("Confirm execution failed")
        
        try:
            response = response_ptr.contents
            
            if response.error_message:
                error_msg = response.error_message.decode('utf-8')
                raise RuntimeError(f"Confirm failed: {error_msg}")
            
            return response.success == 1
        
        finally:
            self._lib.ee_response_free(response_ptr)
