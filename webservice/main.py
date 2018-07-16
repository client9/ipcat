import webapp2
import logging
from IpDatabase import IpDatabase

class ipCat(webapp2.RequestHandler):
    def __init__(self, request, response):
        self.initialize(request, response)
        logging.debug("Instance of ipcat handler created")

    def get(self):
        if self.app.db.needs_update(self.request.get("cachelimit", 86400)):
            self.app.db.update()
        hit = self.app.db.find(self.request.get("ip", "1.1.1.1"))
        self.response.headers['Content-Type'] = 'text/plain'
        self.response.write(hit)

class application(webapp2.WSGIApplication):
    def __init__(self, routes, debug=False, config=None):
        webapp2.WSGIApplication.__init__(self, routes, debug, config)
    	self.db = IpDatabase().generate()


app = application([
    ('/ipCat', ipCat),
    ], debug=True)
