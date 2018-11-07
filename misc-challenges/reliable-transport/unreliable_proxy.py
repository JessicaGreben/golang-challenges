# -*- coding: utf-8 -*-
import ast
import random
import socket
import sys
import traceback


DROP_RATE = 0.3
CORRUPTION_RATE = 0.5


if __name__ == '__main__':
    if len(sys.argv) != 2:
        print('Usage: python3 unreliable_proxy.py [dest]')
        sys.exit(1)

    dest = int(sys.argv[1])
    print('Anything received will be forwarded to 0.0.0.0:{}'.format(dest))
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        sock.bind(('', 0))
        print('Listening on {}:{}'.format(*sock.getsockname()))

        while True:
            payload, address = sock.recvfrom(4096)
            print('\n⭐⭐ New Packet ⭐⭐\n')
            print(payload)
            
            if random.random() < DROP_RATE:
                print('\n🛑 PACKET DROPPED!\n')
                continue

            if random.random() < CORRUPTION_RATE:
                print('\n🔥 CORRUPTION!\n')
                payload = list(payload)
                payload[random.randrange(len(payload))] ^= 0xff
                payload = bytes(payload)

            if 1337 == address[1]:
                continue
            else:
                print('\n👌 OK\n')

            sock.sendto(payload, ('', dest))
    finally:
        sock.close()

