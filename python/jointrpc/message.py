from typing import Any, List, Dict, Union
import abc

class Message:
    @abc.abstractmethod
    def encode(self) -> Dict[str, Any]:
        raise NotImplementedError

class Request(Message):
    id: Any
    method: str
    params: List[Any]
    def __init__(self, id: Any, method: str, *params: Any):
        self.id = id
        self.method = method
        self.params = params

    def encode(self) -> Dict[str, Any]:
        return {
            'version': '2.0',
            'id': self.id,
            'method': self.method,
            'params': self.params
        }

class Notify(Message):
    method: str
    params: List[Any]
    def __init__(self, method: str, *params: Any):
        self.method = method
        self.params = params

    def encode(self) -> Dict[str, Any]:
        return {
            'version': '2.0',
            'method': self.method,
            'params': self.params
        }

class Result(Message):
    id: Any
    result: Any
    def __init__(self, id: Any, result: Any):
        self.id = id
        self.result = result

    def encode(self) -> Dict[str, Any]:
        return {
            'version': '2.0',
            'id': self.id,
            'result': self.result
        }

class Error(Message):
    id: Any
    error: Any

    @classmethod
    def method_not_found(cls, reqid: Any, reason: str=None) -> 'Error':
        errbody = {'code': 404,
                   'reason': f'method not found, {reason}',
                   'retryable': False}
        return Error(reqid, errbody)

    @classmethod
    def server_error(cls, reqid: Any, reason: str=None, retryable: bool=False) -> 'Error':
        errbody = {'code': 500,
                   'reason': f'server error, {reason}',
                   'retryable': retryable}
        return Error(reqid, errbody)

    def __init__(self, id: Any, error: Any):
        assert id
        self.id = id
        self.error = error

    def encode(self) -> Dict[str, Any]:
        return {
            'version': '2.0',
            'id': self.id,
            'error': self.error
        }

class RPCError(Exception):
    code: int
    reason: str
    retryable: bool
    def __init__(self, code: int, reason: str, retryable: bool=False):
        self.code = code
        self.reason = reason
        self.retryable = retryable

    def body(self) -> Dict[str, Any]:
        return {
            'code': self.code,
            'reason': self.reason,
            'retryable': self.retryable
        }

    def error_message(self, reqid: Any) -> 'Error':
        return Error(reqid, self.body())

def parse(payload: Dict[str, Any]) -> 'Message':
    if 'id' in payload:
        if 'method' in payload:
            return Request(payload['id'],
                           payload['method'],
                           *payload.get('params', []))
        elif 'result' in payload:
            return Result(payload['id'], payload['result'])
        elif 'error' in payload:
            return Error(payload['id'], payload['error'])
        else:
            pass
    elif 'method' in payload:
        return Notify(payload['method'],
                      *payload.get('params', []))
    raise ValueError('fail to parse jsonrpc message')
