const express = require('express')
const router = express.Router()
const {
  getUserAll,
  create,
  getUserById,
  updateUserById,
  deleteUserById
} = require('../controllers/userController')
// const { v4: uuidv4 } = require('uuid')
// uuidv4()

/* GET users listing. */
router.get('/', function(req, res, next) {
  res.send('respond with a resource')
})

router.get('/all', getUserAll)
router.post('/create', create)
router.get('/:id', getUserById)
router.put('/:id', updateUserById)
router.delete('/:id', deleteUserById)

module.exports = router
