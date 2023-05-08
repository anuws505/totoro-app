const express = require('express')
const router = express.Router()
// const { v4: uuidv4 } = require('uuid')
// uuidv4()

/* GET users listing. */
router.get('/', function(req, res, next) {
  res.send('respond with a resource')
})

module.exports = router
