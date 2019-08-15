from random import choice

from twisted.web import server, resource
from twisted.internet import reactor

import time

class Uninspiring(Exception):
    pass


class Magnificent(resource.Resource):
    isLeaf = True

    def render_GET(self, request):
        if choice([True, True, True, False]):
            if choice([True, True, False]):
                print('Magnificent')
                return "Magnificent!".encode('utf-8')  # Twisted expects bytes.
            else:
                print('Something else')
                return "Something else".encode('utf-8')  # Twisted expects bytes.
        else:
            raise Uninspiring()


class run():
    site = server.Site(Magnificent())
    # https://www.youtube.com/watch?v=a6iW-8xPw3k
    reactor.listenTCP(12345, site)
    reactor.run()


if __name__ == "__main__":
    run()
