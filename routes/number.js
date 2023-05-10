const express = require('express')
const router = express.Router()
const {
  getNumberAll,
  create
} = require('../controllers/numberController')

router.get('/', function(req, res, next) {
  res.send('what?')
})

router.get('/all', getNumberAll)
router.post('/create', create)

module.exports = router
