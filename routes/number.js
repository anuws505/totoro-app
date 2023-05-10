const express = require('express')
const router = express.Router()
const {
  getNumberAll,
  create,
  getNumberById,
  updateNumberById,
  deleteNumberById
} = require('../controllers/numberController')

router.get('/', function(req, res, next) {
  res.send('what?')
})

router.get('/all', getNumberAll)
router.post('/create', create)
router.get('/:id', getNumberById)
router.put('/:id', updateNumberById)
router.delete('/:id', deleteNumberById)

module.exports = router
