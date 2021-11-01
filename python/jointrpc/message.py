from typing import Any, List, Dict, Union
import json
import abc

IdType = Union[str, int]

class Message:
    traceid: str = ''

    @abc.abstractmethod
    def encode(self) -> Dict[str, Any]:
        raise NotImplementedError

    def __str__(self) -> str:
        return json.dumps(self.encode())

class Request(Message):
    id: IdType
    method: str
    params: List[Any]

    def __init__(self, id: Any, method: str, *params: Any, traceid:str=''):
        self.id = id
        self.method = method
        self.params = params
        self.traceid = traceid

    def encode(self) -> Dict[str, Any]:
        data = {
            'version': '2.0',
            'id': self.id,
            'method': self.method,
            'params': self.params
        }
        if self.traceid:
            data['traceid'] = self.traceid
        return data

class Notify(Message):
    method: str
    params: List[Any]

    def __init__(self, method: str, *params: Any, traceid:str=''):
        self.method = method
        self.params = params
        self.traceid = traceid

    def encode(self) -> Dict[str, Any]:
        data = {
            'version': '2.0',
            'method': self.method,
            'params': self.params
        }
        if self.traceid:
            data['traceid'] = self.traceid
        return data

class Result(Message):
    id: IdType
    body: Any

    def __init__(self, id: Any, body: Any, traceid:str=''):
        self.id = id
        self.body = body
        self.traceid = traceid

    @property
    def result(self) -> Any:
        return self.body

    def encode(self) -> Dict[str, Any]:
        data = {
            'version': '2.0',
            'id': self.id,
            'result': self.body
        }
        if self.traceid:
            data['traceid'] = self.traceid
        return data

class Error(Message):
    id: IdType
    code: int
    message: str
    data: Any

    def __init__(self, id: Any, code: int, message: str, data: Any=None, traceid:str=''):
        assert id
        self.id = id
        self.code = code
        self.message = message
        self.data = data
        self.traceid = traceid

    @property
    def error(self) -> Any:
        return {
            'code': self.code,
            'message': self.message,
            'data': self.data,
        }

    def encode(self) -> Dict[str, Any]:
        data = {
            'version': '2.0',
            'id': self.id,
            'error': self.error
        }
        if self.traceid:
            data['traceid'] = self.traceid
        return data

class RPCError(Exception):
    code: int
    message: str
    data: Any

    def __init__(self, code: int, message: str, data: Any=None):
        self.code = code
        self.message = message
        self.data = data

    def body(self) -> Dict[str, Any]:
        body = {
            'code': self.code,
            'message': self.message,
        }
        if self.data is not None:
            body['data'] = self.data
        return body

    def error_message(self, reqid: Any, traceid:str='') -> 'Error':
        return Error(reqid, self.body(), traceid=traceid)

def parse(payload: Dict[str, Any]) -> 'Message':
    traceid = payload.get('traceid', '')

    if 'id' in payload:
        if 'method' in payload:
            params = payload.get('params', [])
            if not isinstance(params, list):
                params = [params]

            return Request(payload['id'],
                           payload['method'],
                           *params,
                           traceid=traceid)
        elif 'result' in payload:
            return Result(payload['id'],
                          payload['result'],
                          traceid=traceid)
        elif 'error' in payload:
            body = payload['error']
            return Error(payload['id'],
                         body['code'],
                         body['message'],
                         data=body.get('data'),
                         traceid=traceid)
        else:
            pass
    elif 'method' in payload:
        params = payload.get('params', [])
        if not isinstance(params, list):
            params = [params]
        return Notify(payload['method'],
                      *params,
                      traceid=traceid)

    raise ValueError('fail to parse jsonrpc message')
