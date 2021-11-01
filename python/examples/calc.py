import logging
import logging.config
import sys
import asyncio


LOGGING = {
    'version': 1,
    'handlers': {
        'console': {
            'class': 'logging.StreamHandler',
            'formatter': 'simple',
            'stream': sys.stdout,
        },
    },
    'formatters':{
        'simple':{
            'format': '[%(asctime)s] %(name)s [%(levelname)s] %(message)s',
        },
    },
    'root': {
        'handlers': ['console'],
        'level': 'DEBUG'
    },
}


add2num_schema = '''
{
   "type": "method",
   "params": ["number", "number"]
}
'''
async def add2num(msg, a, b):
    print('add2num called', a, b)
    return a + b

async def main():
    logging.config.dictConfig(LOGGING)
    from jointrpc.aio import Client
    c = Client(sys.argv[1])
    c.on('calc_add2num', add2num).help('nice').schema(add2num_schema)

    livestream = c.live_stream()
    async def live_ready():
        print("live stream is ready")
    livestream.set_ready_cb(live_ready)
    await livestream.handle()
    print('stop')

if __name__ == '__main__':
    asyncio.run(main())
