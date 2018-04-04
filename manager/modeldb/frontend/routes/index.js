var express = require('express');
var router = express.Router();
var api = require('../util/api.js');

/* GET home page. */
router.get("/", function(req, res, next) {
  res.redirect(process.env.ROOT_PATH + '/projects');
});

module.exports = router;
